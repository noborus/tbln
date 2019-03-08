package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

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

func decrypt(key []byte, cipherText []byte) ([]byte, error) {
	padkey := pad(key)
	block, err := aes.NewCipher(padkey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := cipherText[:aesgcm.NonceSize()]
	plaintextBytes, err := aesgcm.Open(nil, nonce, cipherText[aesgcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
	return plaintextBytes, nil
}

func decryptPrompt(keyFile string) ([]byte, error) {
	prFile, _ := os.Open(keyFile)
	defer prFile.Close()
	data, err := ioutil.ReadAll(prFile)
	if err != nil {
		return nil, err
	}
	privateKey, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	password, err := readPasswordPrompt("password: ")
	if err != nil {
		return nil, err
	}
	priv, err := decrypt(password, []byte(privateKey))
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

func getPublicKey(pubFileName string) ([]byte, error) {
	pubFile, err := os.Open(pubFileName)
	if err != nil {
		return nil, err
	}
	defer pubFile.Close()
	data, err := ioutil.ReadAll(pubFile)
	if err != nil {
		return nil, err
	}
	pub, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return pub, nil
}
