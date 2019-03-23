package key

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
	"golang.org/x/crypto/ssh/terminal"
)

// WritePrivateFile writes a private key.
func WritePrivateFile(fileName string, keyName string, privateKey []byte) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("private key create error: %s", err)
	}
	kt, err := GeneratePrivate(keyName, privateKey)
	if err != nil {
		return fmt.Errorf("generate private key: %s", err)
	}
	err = tbln.WriteAll(file, kt)
	if err != nil {
		return fmt.Errorf("private key write: %s", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("private key close: %s", err)
	}
	return nil
}

// WritePublicFile writes a publickey.
func WritePublicFile(fileName string, keyName string, public []byte) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("public key file create: %s", err)
	}
	pt, err := GeneratePublic(keyName, public)
	pt.Comments = []string{fmt.Sprintf("TBLN Public key")}
	if err != nil {
		return fmt.Errorf("generate public key: %s", err)
	}
	err = tbln.WriteAll(file, pt)
	if err != nil {
		return fmt.Errorf("public key write: %s", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("public key close: %s", err)
	}
	return nil
}

// GetPrivateKey gets the private key.
func GetPrivateKey(privFileName string, keyName string, prompt bool) ([]byte, error) {
	privFile, err := os.Open(privFileName)
	if err != nil {
		return nil, err
	}
	defer privFile.Close()

	t, err := tbln.ReadAll(privFile)
	if err != nil {
		return nil, fmt.Errorf("private key read error %s: %s", privFileName, err)
	}
	var privateKey []byte
	for _, row := range t.Rows {
		if row[0] == keyName {
			privateKey, err = base64.StdEncoding.DecodeString(row[2])
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if privateKey == nil {
		return nil, fmt.Errorf("no matching secret key found: %s", keyName)
	}
	var priv []byte
	if prompt {
		priv, err = decryptPrompt(privateKey)
	} else {
		priv, err = decrypt([]byte(""), privateKey)
	}
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// GeneratePublic write public key in TBLN file format.
func GeneratePublic(keyName string, pubkey []byte) (*tbln.Tbln, error) {
	var err error
	t := tbln.NewTbln()
	err = t.SetNames([]string{"keyname", "algorithm", "publickey"})
	if err != nil {
		return nil, err
	}
	err = t.SetTypes([]string{"text", "text", "text"})
	if err != nil {
		return nil, err
	}
	pEnc := base64.StdEncoding.EncodeToString(pubkey)
	err = t.AddRows([]string{keyName, tbln.ED25519, pEnc})
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GeneratePrivate write private key in TBLN file format.
func GeneratePrivate(keyName string, privkey []byte) (*tbln.Tbln, error) {
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
