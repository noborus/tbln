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
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh/terminal"
)

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "generate a new key pair",
	Long:  `generate a new key pair`,
	Run: func(cmd *cobra.Command, args []string) {
		genkey(cmd, args)
	},
}

func init() {
	genkeyCmd.PersistentFlags().StringP("name", "n", "", "Key Name")
	rootCmd.AddCommand(genkeyCmd)
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

func genkey(cmd *cobra.Command, args []string) error {
	var name string
	var err error
	if len(args) > 0 {
		name = args[0]
	}
	if name, err = cmd.PersistentFlags().GetString("name"); err != nil {
		log.Fatal(err)
	}
	if name == "" {
		return fmt.Errorf("should be name")
	}
	fmt.Printf("password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Printf("\n")
	fmt.Printf("confirm: ")
	confirm, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	if string(password) != string(confirm) {
		return fmt.Errorf("password invalid")
	}
	fmt.Printf("\n")
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	cipherText, err := encrypt(password, private)
	if err != nil {
		return err
	}
	privateFile, err := os.OpenFile(name+".key", os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewWriter(privateFile)
	ps := base64.StdEncoding.EncodeToString(cipherText)
	_, err = buf.WriteString(ps)
	if err != nil {
		log.Fatal(err)
	}
	buf.Flush()
	privateFile.Close()

	pubFile, err := os.OpenFile(name+".pub", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	buf2 := bufio.NewWriter(pubFile)
	ps2 := base64.StdEncoding.EncodeToString([]byte(public))
	_, err = buf2.WriteString(ps2)
	if err != nil {
		log.Fatal(err)
	}
	buf2.Flush()
	pubFile.Close()
	return nil
}
