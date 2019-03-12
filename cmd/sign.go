package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:          "sign [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Sign a TBLN file with a private key.",
	Long: `Sign a TBLN file with a private key.
Sign the hash value with the ED25519 private key.
Generate hash first, if there is no hash value yet(default is SHA256).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return signFile(cmd, args)
	},
}

func init() {
	signCmd.PersistentFlags().StringSliceP("hash", "a", []string{}, "Hash algorithm(sha256 or sha512)")
	signCmd.PersistentFlags().StringP("key", "k", "", "Key File")
	signCmd.PersistentFlags().BoolP("quiet", "q", false, "Do not prompt for password.")
	signCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(signCmd)
}

func signFile(cmd *cobra.Command, args []string) error {
	var err error
	var fileName, keyFile string
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
	if keyFile, err = cmd.PersistentFlags().GetString("key"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if keyFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be keyFile")
	}
	var quiet bool
	if quiet, err = cmd.PersistentFlags().GetBool("quiet"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	var privKey []byte
	if quiet {
		privKey, err = decrypt([]byte(""), keyFile)
	} else {
		privKey, err = decryptPrompt(keyFile)
	}
	if err != nil {
		return err
	}
	keyName := filepath.Base(keyFile[:len(keyFile)-len(filepath.Ext(keyFile))])

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
		at.SumHash(hash)
	}
	at.Sign(keyName, privKey)
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}

	return nil
}
