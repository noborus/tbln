package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"
)

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey [Key Name]",
	Short: "generate a new key pair",
	Long: `Generate a new Ed25519 keypair.
Write privake key(Key Name+".key") and public key(Key Name+".pub") files
based on the Key Name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genkey(cmd, args)
	},
}

func init() {
	genkeyCmd.PersistentFlags().StringP("keyName", "k", "", "Key Name")
	rootCmd.AddCommand(genkeyCmd)
}

func genkey(cmd *cobra.Command, args []string) error {
	var keyName string
	var err error
	if keyName, err = cmd.PersistentFlags().GetString("keyName"); err != nil {
		return err
	}
	if keyName == "" && len(args) > 0 {
		keyName = args[0]
	}
	if keyName == "" {
		return fmt.Errorf("must be Key Name")
	}

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

	pubKeyName := keyName + ".pub"
	priKeyName := keyName + ".key"

	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	cipherText, err := encrypt(password, private)
	if err != nil {
		return err
	}

	privateFile, err := os.OpenFile(priKeyName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	buf := bufio.NewWriter(privateFile)
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

	pubFile, err := os.OpenFile(pubKeyName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	buf2 := bufio.NewWriter(pubFile)
	_, err = buf2.WriteString(keyName + ":")
	if err != nil {
		return err
	}
	ps2 := base64.StdEncoding.EncodeToString([]byte(public))
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
	fmt.Printf("write %s %s file\n", pubKeyName, priKeyName)
	return nil
}
