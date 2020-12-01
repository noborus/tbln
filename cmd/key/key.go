// Package key writes and reads the private/public key.
//
// The secret key is AES encrypted by password.
// This package does not include the public key read,
// it will be read by the keystore.
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
	"golang.org/x/crypto/ed25519"
	"golang.org/x/term"
)

// GenerateKey generates private key and public key
func GenerateKey() ([]byte, []byte, error) {
	return ed25519.GenerateKey(nil)
}

// WritePrivateFile writes a private key.
func WritePrivateFile(fileName string, keyName string, privateKey []byte) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("private key create error: %s", err)
	}
	kt, err := genTBLNPrivate(keyName, privateKey)
	if err != nil {
		return fmt.Errorf("generate TBLN private: %s", err)
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

// genTBLNPrivate write private key in TBLN file format.
func genTBLNPrivate(keyName string, privkey []byte) (*tbln.TBLN, error) {
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
	t := tbln.NewTBLN()
	t.Comments = []string{"TBLN Private key"}
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

// WritePublicFile writes a publickey.
func WritePublicFile(fileName string, keyName string, public []byte) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("public key file create: %s", err)
	}
	pt, err := GenTBLNPublic(keyName, public)
	if err != nil {
		return fmt.Errorf("generate public key: %s", err)
	}
	pt.Comments = []string{"TBLN Public key"}
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

// GenTBLNPublic write public key in TBLN file format.
func GenTBLNPublic(keyName string, pubkey []byte) (*tbln.TBLN, error) {
	var err error
	if keyName == "" {
		return nil, fmt.Errorf("keyname is required")
	}
	if len(pubkey) == 0 {
		return nil, fmt.Errorf("pubkey is required")
	}
	t := tbln.NewTBLN()
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

// PR determines the interface of the password reader.
var PR PasswordReader = StdInPasswordReader{}

func readPasswordPrompt(prompt string) ([]byte, error) {
	return PR.ReadPasswordPrompt(prompt)
}

// PasswordReader is an interface for reading a password.
type PasswordReader interface {
	ReadPasswordPrompt(prompt string) ([]byte, error)
}

// StdInPasswordReader is a standard password reader.
type StdInPasswordReader struct {
}

// ReadPasswordPrompt prints a password prompt and returns the password.
func (pr StdInPasswordReader) ReadPasswordPrompt(prompt string) ([]byte, error) {
	fd := int(os.Stdin.Fd())
	fmt.Fprint(os.Stderr, prompt)
	password, err := term.ReadPassword(fd)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(os.Stderr, "\n")
	return password, nil
}
