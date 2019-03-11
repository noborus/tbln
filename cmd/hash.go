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
	Long:         `Add hash value to TBLN file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return addHash(cmd, args)
	},
}

func init() {
	hashCmd.PersistentFlags().BoolP("sha256", "", false, "SHA256 Hash")
	hashCmd.PersistentFlags().BoolP("sha512", "", false, "SHA512 Hash")
	hashCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(hashCmd)
}

func addHash(cmd *cobra.Command, args []string) error {
	var fileName string
	var sha256, sha512 bool
	var err error
	if sha256, err = cmd.PersistentFlags().GetBool("sha256"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if sha512, err = cmd.PersistentFlags().GetBool("sha512"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if sha256 == false && sha512 == false {
		sha256 = true
	}
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
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}

	return nil
}
