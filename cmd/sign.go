package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a TBLN file with a private key",
	Long:  `Sign a TBLN file with a private key`,
	Run: func(cmd *cobra.Command, args []string) {
		signFile(cmd, args)
	},
}

func init() {
	signCmd.PersistentFlags().StringP("key", "k", "", "Key File")
	signCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")

	rootCmd.AddCommand(signCmd)
}

func signFile(cmd *cobra.Command, args []string) {
	var fileName, keyFile string
	var password []byte
	var err error
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		log.Fatal(err)
	}
	if len(args) > 0 {
		fileName = args[0]
	}
	if keyFile, err = cmd.PersistentFlags().GetString("key"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("password: ")
	password, err = terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	prFile, _ := os.Open(keyFile)
	defer prFile.Close()
	fi, _ := prFile.Stat()
	size := fi.Size()
	data := make([]byte, size)
	prFile.Read(data)
	privateKey, _ := base64.StdEncoding.DecodeString(string(data))
	priv, err := decrypt(password, []byte(privateKey))
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = at.SumHash(tbln.SHA256)
	if err != nil {
		log.Fatal(err)
	}
	at.Sign(priv)
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}

	return
}
