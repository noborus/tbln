package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:          "verify [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Verify signature and checksum of TBLN file",
	Long:         `Verify signature and checksum of TBLN file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return verify(cmd, args)
	},
}

func init() {
	verifyCmd.PersistentFlags().BoolP("force-verify-sign", "", false, "Force signature verification")
	verifyCmd.PersistentFlags().BoolP("no-verify-sign", "", false, "ignore signature verify")
	verifyCmd.PersistentFlags().StringP("output", "o", "", "Write to file instead of stdout")
	verifyCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(verifyCmd)
}

func verify(cmd *cobra.Command, args []string) error {
	_, err := verifiedTbln(cmd, args)
	return err
}

func verifiedTbln(cmd *cobra.Command, args []string) (*tbln.Tbln, error) {
	var err error
	var fileName string
	var forcesign, nosign bool
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		cmd.SilenceUsage = false
		return nil, err
	}
	if fileName == "" && len(args) > 0 {
		fileName = args[0]
	}
	if fileName == "" {
		cmd.SilenceUsage = false
		return nil, fmt.Errorf("require filename")
	}
	if nosign, err = cmd.PersistentFlags().GetBool("no-verify-sign"); err != nil {
		return nil, err
	}
	if forcesign, err = cmd.PersistentFlags().GetBool("force-verify-sign"); err != nil {
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	if at.TableName() == "" {
		at.SetTableName(filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))]))
	}
	if forcesign || (!nosign && (len(at.Signs) != 0)) {
		pub, err := getPublicKey(pubfile, keyname)
		if err != nil {
			return nil, err
		}
		if !at.VerifySignature(keyname, pub) {
			return nil, fmt.Errorf("signature verification failure")
		}
		fmt.Fprintln(os.Stderr, "Signature verification successful")
	} else {
		if !at.Verify() {
			return nil, fmt.Errorf("verification failure")
		}
		fmt.Fprintln(os.Stderr, "Verification successful")
	}
	return at, nil
}