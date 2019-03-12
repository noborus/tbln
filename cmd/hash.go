package cmd

import (
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:          "hash [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Add hash value to TBLN file",
	Long: `Add hash value to TBLN file.
Select SHA256, SHA512, or both HASH values to the TBLN file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return addHash(cmd, args)
	},
}

func init() {
	hashCmd.PersistentFlags().StringSliceP("hash", "a", []string{"sha256"}, "Hash algorithm(sha256 or sha512)")
	hashCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(hashCmd)
}

func addHash(cmd *cobra.Command, args []string) error {
	var err error
	var hashes []string
	if hashes, err = cmd.PersistentFlags().GetStringSlice("hash"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	var fileName string
	if len(args) > 0 {
		fileName = args[0]
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
	for _, hash := range hashes {
		err = at.SumHash(hash)
		if err != nil {
			return err
		}
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}

	return nil
}
