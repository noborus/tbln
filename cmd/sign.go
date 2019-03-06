package cmd

import (
	"fmt"
	"os"

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
	signCmd.PersistentFlags().StringP("key", "k", "", "Key File")
	signCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(signCmd)
}

func signFile(cmd *cobra.Command, args []string) error {
	var fileName, keyFile string
	var err error
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

	priv, err := decryptPrompt(keyFile)
	if err != nil {
		return err
	}
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
	err = at.SumHash(tbln.SHA256)
	if err != nil {
		return err
	}
	at.Sign(priv)
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}

	return nil
}
