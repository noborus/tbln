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
	hashCmd.PersistentFlags().StringSliceP("enable-target", "", nil, "Hash target extra item (all or each name)")
	hashCmd.PersistentFlags().StringSliceP("disable-target", "", nil, "Hash extra items not to be targeted (all or each name)")
	hashCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(hashCmd)
}

func hash(at *tbln.Tbln, cmd *cobra.Command, args []string) error {
	var err error
	var hashes []string
	if hashes, err = cmd.PersistentFlags().GetStringSlice("hash"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	var enableTarget []string
	if enableTarget, err = cmd.PersistentFlags().GetStringSlice("enable-target"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	var disableTarget []string
	if disableTarget, err = cmd.PersistentFlags().GetStringSlice("disable-target"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	for _, v := range enableTarget {
		if v == "all" {
			at.AllTargetHash(true)
			break
		}
		if _, ok := at.Extras[v]; ok {
			at.ToTargetHash(v, true)
		}
	}
	for _, v := range disableTarget {
		if v == "all" {
			at.AllTargetHash(false)
			break
		}
		if _, ok := at.Extras[v]; ok {
			at.ToTargetHash(v, false)
		}
	}

	if len(hashes) == 0 {
		for k := range at.Hashes {
			hashes = append(hashes, k)
		}
	}
	if len(hashes) == 0 {
		hashes = append(hashes, "sha256")
	}
	for _, hash := range hashes {
		err = at.SumHash(hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func addHash(cmd *cobra.Command, args []string) error {
	var err error

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
	hash(at, cmd, args)
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}

	return nil
}
