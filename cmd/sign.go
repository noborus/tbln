package cmd

import (
	"fmt"

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
	RunE: sign,
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

func sign(cmd *cobra.Command, args []string) error {
	var err error
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

	at, err := readTbln(fileName, cmd)
	if err != nil {
		return err
	}
	at, err = hashFile(at, cmd)
	if err != nil {
		return err
	}
	at, err = signFile(at, cmd)
	if err != nil {
		return err
	}
	return outputFile(at, cmd)
}

func signFile(at *tbln.Tbln, cmd *cobra.Command) (*tbln.Tbln, error) {
	var err error
	var quiet bool
	if quiet, err = cmd.PersistentFlags().GetBool("quiet"); err != nil {
		cmd.SilenceUsage = false
		return nil, err
	}
	privKey, err := key.GetPrivateKey(SecFile, KeyName, !quiet)
	if err != nil {
		return nil, err
	}
	_, err = at.Sign(KeyName, privKey)
	if err != nil {
		return nil, err
	}
	return at, nil
}
