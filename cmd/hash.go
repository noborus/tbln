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
	RunE: addHash,
}

func init() {
	hashCmd.PersistentFlags().StringSliceP("hash", "a", []string{"sha256"}, "hash algorithm(sha256 or sha512)")
	hashCmd.PersistentFlags().StringSliceP("enable-target", "", nil, "hash target extra item (all or each name)")
	hashCmd.PersistentFlags().StringSliceP("disable-target", "", nil, "hash extra items not to be targeted (all or each name)")
	hashCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	hashCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")

	rootCmd.AddCommand(hashCmd)
}

func hash(at *tbln.Tbln, cmd *cobra.Command) error {
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
	var output string
	if output, err = cmd.PersistentFlags().GetString("output"); err != nil {
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
	err = hash(at, cmd)
	if err != nil {
		return err
	}
	var out *os.File
	if output == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(output)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return tbln.WriteAll(out, at)
}
