package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/noborus/tbln"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh/terminal"
)

func generateKey(keyName string, overwrite bool) error {
	if _, err := os.Stat(KeyPath); os.IsNotExist(err) {
		log.Printf("mkdir [%s]", KeyPath)
		err := os.Mkdir(KeyPath, 0700)
		if err != nil {
			return err
		}
	}
	_, err := os.Stat(SecKey)
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("%s file already exists", SecKey)
	}
	_, err = os.Stat(PubFile)
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("%s file already exists", PubFile)
	}

	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	privateFile, err := os.OpenFile(SecKey, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("private key create error: %s: %s", SecKey, err)
	}
	kt, err := generatePrivate(keyName, private)
	if err != nil {
		return fmt.Errorf("generate private error: %s: %s", SecKey, err)
	}
	err = tbln.WriteAll(privateFile, kt)
	if err != nil {
		return fmt.Errorf("private key write error: %s: %s", SecKey, err)
	}
	err = privateFile.Close()
	if err != nil {
		return fmt.Errorf("private key write error: %s: %s", SecKey, err)
	}
	log.Printf("write %s file\n", SecKey)

	pub, err := os.OpenFile(PubFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	pt, err := generatePublic(keyName, public)
	if err != nil {
		return fmt.Errorf("generate key create: %s: %s", PubFile, err)
	}
	err = tbln.WriteAll(pub, pt)
	if err != nil {
		return fmt.Errorf("public key write error: %s: %s", PubFile, err)
	}
	err = pub.Close()
	if err != nil {
		return fmt.Errorf("public key write error: %s: %s", PubFile, err)
	}

	/* regist public keys */
	pubk, err := os.OpenFile(PubKeys, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer pubk.Close()

	pkeys, err := tbln.ReadAll(pubk)
	if err != nil {
		return fmt.Errorf("public keys read error %s: %s", PubKeys, err)
	}
	pkeys, err = registPublicKey(pkeys, pt)
	if err != nil {
		return fmt.Errorf("resist public keys %s: %s", PubKeys, err)
	}
	_, err = pubk.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("public keys write error: %s: %s", PubKeys, err)
	}
	err = tbln.WriteAll(pubk, pkeys)
	if err != nil {
		return fmt.Errorf("public keys write error: %s: %s", PubKeys, err)
	}

	log.Printf("write %s file\n", PubKeys)
	return nil
}

func generatePrivate(keyName string, privkey []byte) (*tbln.Tbln, error) {
	password, err := readPasswordPrompt("password: ")
	if err != nil {
		return nil, err
	}
	confirm, err := readPasswordPrompt("confirm: ")
	if err != nil {
		return nil, err
	}
	if string(password) != string(confirm) {
		return nil, fmt.Errorf("password invalid")
	}

	cipherText, err := encrypt(password, privkey)
	if err != nil {
		return nil, err
	}

	t := tbln.NewTbln()
	t.Comments = []string{fmt.Sprintf("TBLN Pvivate key")}
	psEnc := base64.StdEncoding.EncodeToString(cipherText)
	err = t.SetNames([]string{"keyname", "algorithm", "privatekey"})
	if err != nil {
		return nil, err
	}
	err = t.SetTypes([]string{"text", "text", "text"})
	if err != nil {
		return nil, err
	}
	err = t.AddRows([]string{keyName, tbln.ED25519, psEnc})
	if err != nil {
		return nil, err
	}

	return t, nil
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

func getPrivateKeyFile(privFileName string, keyName string) ([]byte, error) {
	privFile, err := os.Open(privFileName)
	if err != nil {
		return nil, err
	}
	defer privFile.Close()

	t, err := tbln.ReadAll(privFile)
	if err != nil {
		return nil, fmt.Errorf("private key read error %s: %s", privFileName, err)
	}
	for _, row := range t.Rows {
		if row[0] == keyName {
			privateKey, err := base64.StdEncoding.DecodeString(row[2])
			if err != nil {
				return nil, err
			}
			return privateKey, nil
		}
	}
	return nil, fmt.Errorf("no matching secret key found: %s", keyName)
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
		return nil, fmt.Errorf("read password :%s", err)
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

func generatePublic(keyName string, pubkey []byte) (*tbln.Tbln, error) {
	var err error
	t := tbln.NewTbln()
	t.Comments = []string{fmt.Sprintf("TBLN Public key")}
	pubEnc := base64.StdEncoding.EncodeToString(pubkey)
	err = t.SetNames([]string{"keyname", "algorithm", "publickey"})
	if err != nil {
		return nil, err
	}
	err = t.SetTypes([]string{"text", "text", "text"})
	if err != nil {
		return nil, err
	}
	err = t.AddRows([]string{keyName, tbln.ED25519, pubEnc})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func registPublicKey(pkeys *tbln.Tbln, addkey *tbln.Tbln) (*tbln.Tbln, error) {
	pkeys.Comments = []string{fmt.Sprintf("TBLN Public keys")}
	for _, row := range addkey.Rows {
		pkeys.AddRows(row)
	}
	return pkeys, nil
}

func getPublicKey(keyName string) ([]byte, error) {
	pubFile, err := os.Open(PubKeys)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", PubKeys, err)
	}
	defer pubFile.Close()
	pt, err := tbln.ReadAll(pubFile)
	if err != nil {
		return nil, fmt.Errorf("public keys read error %s: %s", PubKeys, err)
	}
	for _, row := range pt.Rows {
		if row[0] == keyName {
			return base64.StdEncoding.DecodeString(row[2])
		}
	}
	return nil, fmt.Errorf("no matching public key found: %s", keyName)
}
