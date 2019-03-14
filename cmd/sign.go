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
	Long: `Sign a TBLN file with a private key.
Sign the hash value with the ED25519 private key.
Generate hash first, if there is no hash value yet(default is SHA256).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return signFile(cmd, args)
	},
}

func init() {
	signCmd.PersistentFlags().StringSliceP("hash", "a", []string{}, "Hash algorithm(sha256 or sha512)")
	signCmd.PersistentFlags().BoolP("quiet", "q", false, "Do not prompt for password.")
	signCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(signCmd)
}

func signFile(cmd *cobra.Command, args []string) error {
	var err error
	var fileName string
	var hashes []string
	if hashes, err = cmd.PersistentFlags().GetStringSlice("hash"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if len(args) > 0 {
		fileName = args[0]
	}
	if seckey == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be secret key File")
	}
	var quiet bool
	if quiet, err = cmd.PersistentFlags().GetBool("quiet"); err != nil {
		cmd.SilenceUsage = false
		return err
	}

	privKey, err := getPrivateKeyFile(seckey, keyname)
	if err != nil {
		return err
	}
	var priv []byte
	if quiet {
		priv, err = decrypt([]byte(""), privKey)
	} else {
		priv, err = decryptPrompt(privKey)
	}
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
	at.Sign(keyname, priv)
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}

	return nil
}
