package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// genkeyCmd represents the generate key pair command
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

	err = generateKey(keyName)
	if err != nil {
		return err
	}
	return nil
}
