package cmd

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

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
	buf := bufio.NewWriter(privateFile)
	cipherText, err := encrypt(password, privkey)
	if err != nil {
		return err
	}
	ps := base64.StdEncoding.EncodeToString(cipherText)
	_, err = buf.WriteString(ps)
	if err != nil {
		return err
	}
	_, err = buf.WriteString("\n")
	if err != nil {
		return err
	}
	err = buf.Flush()
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

func decrypt(key []byte, keyFile string) ([]byte, error) {
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
	padkey := pad(key)
	block, err := aes.NewCipher(padkey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := privateKey[:aesgcm.NonceSize()]
	plaintextBytes, err := aesgcm.Open(nil, nonce, privateKey[aesgcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
	return plaintextBytes, nil
}

func decryptPrompt(keyFile string) ([]byte, error) {
	password, err := readPasswordPrompt("password: ")
	if err != nil {
		return nil, err
	}
	priv, err := decrypt(password, keyFile)
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

func getPublicKey(pubFileName string, keyName string) ([]byte, error) {
	if len(pubFileName) == 0 {
		return nil, fmt.Errorf("must be public key file")
	}
	pubs, err := readPublicKeyFile(pubFileName)
	if err != nil {
		return nil, err
	}
	if pub, ok := pubs[keyName]; ok {
		return pub, nil
	}
	return nil, fmt.Errorf("not public key %s", keyName)
}

func writePublicKeyFile(pubFileName string, keyName string, pubkey []byte) error {
	pubFile, err := os.OpenFile(keyName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	buf2 := bufio.NewWriter(pubFile)
	_, err = buf2.WriteString(keyName + ":")
	if err != nil {
		return err
	}
	ps2 := base64.StdEncoding.EncodeToString([]byte(pubkey))
	_, err = buf2.WriteString(ps2)
	if err != nil {
		return err
	}
	err = buf2.Flush()
	if err != nil {
		return err
	}

	err = pubFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func readPublicKeyFile(pubFileName string) (map[string][]byte, error) {
	pubFile, err := os.Open(pubFileName)
	if err != nil {
		return nil, err
	}
	defer pubFile.Close()
	pubs := make(map[string][]byte)
	scanner := bufio.NewScanner(pubFile)
	for scanner.Scan() {
		txt := scanner.Text()
		if txt == "" {
			continue
		}
		data := strings.SplitN(txt, ":", 2)
		if len(data) != 2 {
			return nil, fmt.Errorf("public key parser error")
		}
		name := data[0]
		pubEnc := data[1]
		pubkey, err := base64.StdEncoding.DecodeString(string(pubEnc))
		pubs[name] = pubkey
		if err != nil {
			return nil, err
		}
	}
	return pubs, nil
}
