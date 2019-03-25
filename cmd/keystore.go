package cmd

import (
	"fmt"
	"strconv"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/cmd/keystore"
	"github.com/spf13/cobra"
)

// keystoreCmd represents the keystore command
var keystoreCmd = &cobra.Command{
	Use:   "keystore",
	Short: "keystore is a command to operate keystore",
	Long:  `keystore is a command to operate keystore`,
}

func init() {
	rootCmd.AddCommand(keystoreCmd)
	keystoreCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List of public keys registered in keystore.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list()
		},
	})
	keystoreCmd.AddCommand(&cobra.Command{
		Use:   "add [public_key.pub]",
		Args:  cobra.ExactArgs(1),
		Short: "Add public key to keystore.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := keystore.AddKey(KeyStore, args[0])
			if err != nil {
				return err
			}
			return list()
		},
	})
	keystoreCmd.AddCommand(&cobra.Command{
		Use:   "del [keyName] <num>",
		Short: "Delete public key from keystore.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) == 0 {
				return fmt.Errorf("requires delete key name")
			}
			fileName := args[0]
			num := -1
			if len(args) > 1 {
				if num, err = strconv.Atoi(args[1]); err != nil {
					return err
				}
			}
			err = keystore.DelKey(KeyStore, fileName, num)
			if err != nil {
				return err
			}
			return list()
		},
	})
}

func list() error {
	list, err := keystore.List(KeyStore)
	if err != nil {
		return err
	}
	for _, row := range list {
		fmt.Println(tbln.JoinRow(row))
	}
	return nil
}
