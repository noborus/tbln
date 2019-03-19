package cmd

import (
	"github.com/spf13/cobra"
)

// genkeyCmd represents the generate key pair command
var genkeyCmd = &cobra.Command{
	Use:          "genkey [Key Name]",
	SilenceUsage: true,
	Short:        "Generate a new key pair",
	Long: `Generate a new Ed25519 keypair.
Write privake key(Key Name+".key") and public key(Key Name+".pub") files
based on the Key Name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genkey(cmd, args)
	},
}

func init() {
	genkeyCmd.PersistentFlags().BoolP("force", "", false, "Overwrites even if the file already exists.")

	rootCmd.AddCommand(genkeyCmd)
}

func genkey(cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		keyname = args[0]
	}
	var overwrite bool
	if overwrite, err = cmd.PersistentFlags().GetBool("force"); err != nil {
		return err
	}
	err = generateKey(keyname, overwrite)
	if err != nil {
		return err
	}
	return nil
}
