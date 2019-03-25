package cmd

import (
	"fmt"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/cmd/key"

	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:          "sign [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Sign a TBLN file with a private key",
	Long: `Sign a TBLN file with a private key.
Sign the hash value with the ED25519 private key.
Generate hash first, if there is no hash value yet(default is SHA256).`,
	RunE: signFile,
}

func init() {
	signCmd.PersistentFlags().StringSliceP("hash", "a", []string{"sha256"}, "hash algorithm(sha256 or sha512)")
	signCmd.PersistentFlags().StringSliceP("enable-target", "", nil, "hash target extra item (all or each name)")
	signCmd.PersistentFlags().StringSliceP("disable-target", "", nil, "hash extra items not to be targeted (all or each name)")
	signCmd.PersistentFlags().BoolP("quiet", "q", false, "do not prompt for password.")
	signCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	signCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")

	rootCmd.AddCommand(signCmd)
}

func signFile(cmd *cobra.Command, args []string) error {
	var err error
	var output string
	if output, err = cmd.PersistentFlags().GetString("output"); err != nil {
		return err
	}
	var fileName string
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if len(args) > 0 {
		fileName = args[0]
	}
	if fileName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("require filename")
	}
	var quiet bool
	if quiet, err = cmd.PersistentFlags().GetBool("quiet"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	privKey, err := key.GetPrivateKey(SecFile, KeyName, !quiet)
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
	err = hash(at, cmd, args)
	if err != nil {
		return err
	}
	_, err = at.Sign(KeyName, privKey)
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
