package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/noborus/tbln"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh/terminal"
)

func generateKey(keyName string) error {
	pubKeyFileName := keyName + ".pub"
	priKeyFileName := keyName + ".key"
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	err = writePrivateFile(priKeyFileName, keyName, private)
	if err != nil {
		return err
	}
	err = writePublicKeyFile(pubKeyFileName, keyName, public)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "write %s %s file\n", pubKeyFileName, priKeyFileName)
	return nil
}

func writePrivateFile(privFileName string, keyName string, privkey []byte) error {
	password, err := readPasswordPrompt("password: ")
	if err != nil {
		return err
	}
	confirm, err := readPasswordPrompt("confirm: ")
	if err != nil {
		return err
	}
	if string(password) != string(confirm) {
		return fmt.Errorf("password invalid")
	}

	privateFile, err := os.OpenFile(privFileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	cipherText, err := encrypt(password, privkey)
	if err != nil {
		return err
	}

	t := tbln.NewTbln()
	t.Comments = []string{fmt.Sprintf("TBLN Pvivate key")}
	psEnc := base64.StdEncoding.EncodeToString(cipherText)
	err = t.SetNames([]string{"keyname", "alogrithm", "privatekey"})
	if err != nil {
		return err
	}
	err = t.SetTypes([]string{"text", "text", "text"})
	if err != nil {
		return err
	}
	err = t.AddRows([]string{keyName, tbln.ED25519, psEnc})
	if err != nil {
		return err
	}
	err = tbln.WriteAll(privateFile, t)
	if err != nil {
		return err
	}
	err = privateFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func pad(b []byte) []byte {
	padSize := aes.BlockSize - (len(b) % (aes.BlockSize))
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(b, pad...)
}

func encrypt(key []byte, msg []byte) ([]byte, error) {
	padkey := pad(key)
	block, err := aes.NewCipher(padkey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherText := aesgcm.Seal(nil, nonce, msg, nil)
	cipherText = append(nonce, cipherText...)
	return cipherText, nil
}

func getPrivateKeyFile(privKeyFile string, keyName string) ([]byte, error) {
	privFile, err := os.Open(privKeyFile)
	if err != nil {
		return nil, err
	}
	defer privFile.Close()

	t, err := tbln.ReadAll(privFile)
	if err != nil {
		return nil, err
	}
	for _, row := range t.Rows {
		if row[0] == keyName {
			privateKey, err := base64.StdEncoding.DecodeString(string(row[2]))
			if err != nil {
				return nil, err
			}
			return privateKey, nil
		}
	}
	return nil, fmt.Errorf("not private key :%s", keyName)
}

func decrypt(key []byte, privKey []byte) ([]byte, error) {
	padkey := pad(key)
	block, err := aes.NewCipher(padkey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := privKey[:aesgcm.NonceSize()]
	plaintextBytes, err := aesgcm.Open(nil, nonce, privKey[aesgcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
	return plaintextBytes, nil
}

func decryptPrompt(privKey []byte) ([]byte, error) {
	password, err := readPasswordPrompt("password: ")
	if err != nil {
		return nil, err
	}
	priv, err := decrypt(password, privKey)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func readPasswordPrompt(prompt string) ([]byte, error) {
	fd := int(os.Stdin.Fd())
	fmt.Fprint(os.Stderr, prompt)
	password, err := terminal.ReadPassword(fd)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(os.Stderr, "\n")
	return password, nil
}

func writePublicKeyFile(pubFileName string, keyName string, pubkey []byte) error {
	pubFile, err := os.OpenFile(pubFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	t := tbln.NewTbln()
	t.Comments = []string{fmt.Sprintf("TBLN Public key")}
	pubEnc := base64.StdEncoding.EncodeToString([]byte(pubkey))
	err = t.SetNames([]string{"keyname", "alogrithm", "publickey"})
	if err != nil {
		return err
	}
	err = t.SetTypes([]string{"text", "text", "text"})
	if err != nil {
		return err
	}
	err = t.AddRows([]string{keyName, tbln.ED25519, pubEnc})
	if err != nil {
		return err
	}
	err = tbln.WriteAll(pubFile, t)
	if err != nil {
		return err
	}

	err = pubFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func getPublicKey(pubFileName string, keyName string) ([]byte, error) {
	if len(pubFileName) == 0 {
		return nil, fmt.Errorf("must be public key file")
	}
	pubFile, err := os.Open(pubFileName)
	if err != nil {
		return nil, err
	}
	defer pubFile.Close()
	pt, err := tbln.ReadAll(pubFile)
	if err != nil {
		return nil, err
	}
	for _, row := range pt.Rows {
		if row[0] == keyName {
			return base64.StdEncoding.DecodeString(string(row[2]))
		}
	}
	return nil, fmt.Errorf("not public key %s", keyName)
}
