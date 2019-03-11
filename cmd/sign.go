package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:          "sign [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Sign a TBLN file with a private key",
	Long:         `Sign a TBLN file with a private key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return signFile(cmd, args)
	},
}

func init() {
	signCmd.PersistentFlags().BoolP("sha256", "", false, "SHA256 Hash")
	signCmd.PersistentFlags().BoolP("sha512", "", false, "SHA512 Hash")
	signCmd.PersistentFlags().StringP("key", "k", "", "Key File")
	signCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(signCmd)
}

func signFile(cmd *cobra.Command, args []string) error {
	var err error
	var fileName, keyFile string
	var sha256, sha512 bool
	if sha256, err = cmd.PersistentFlags().GetBool("sha256"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if sha512, err = cmd.PersistentFlags().GetBool("sha512"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if len(args) > 0 {
		fileName = args[0]
	}
	if keyFile, err = cmd.PersistentFlags().GetString("key"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if keyFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be keyFile")
	}

	privKey, err := decryptPrompt(keyFile)
	if err != nil {
		return err
	}
	keyName := filepath.Base(keyFile[:len(keyFile)-len(filepath.Ext(keyFile))])

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	for k := range at.Hashes {
		if k == "sha256" {
			sha256 = true
		}
		if k == "sha512" {
			sha512 = true
		}
	}
	if sha256 == false && sha512 == false {
		sha256 = true
	}
	if sha256 {
		err = at.SumHash(tbln.SHA256)
		if err != nil {
			return err
		}
	}
	if sha512 {
		err = at.SumHash(tbln.SHA512)
		if err != nil {
			return err
		}
	}
	at.Sign(keyName, privKey)
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}

	return nil
}
