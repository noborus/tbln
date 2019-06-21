package cmd

import (
	"fmt"
	"log"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/cmd/keystore"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:          "verify [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Verify signature and checksum of TBLN file",
	Long:         `Verify signature and checksum of TBLN file`,
	RunE:         verify,
}

func init() {
	verifyCmd.PersistentFlags().BoolP("force-verify-sign", "", false, "force signature verification")
	verifyCmd.PersistentFlags().BoolP("no-verify-sign", "", false, "ignore signature verification")
	verifyCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(verifyCmd)
}

func verify(cmd *cobra.Command, args []string) error {
	cmd.PersistentFlags().StringP("output", "o", "", "")
	_, err := verifiedTBLN(cmd, args)
	return err
}

func verifyAllSignature(tb *tbln.TBLN) error {
	for keyName := range tb.Signs {
		pub, err := keystore.Search(KeyStore, keyName)
		if err != nil {
			return err
		}
		if !tb.VerifySignature(keyName, pub) {
			return fmt.Errorf("signature verification failure")
		}
	}
	return nil
}

func verifiedTBLN(cmd *cobra.Command, args []string) (*tbln.TBLN, error) {
	var err error
	var forceSign, noSign, noVerify bool
	if noSign, err = cmd.PersistentFlags().GetBool("no-verify-sign"); err != nil {
		return nil, err
	}
	if forceSign, err = cmd.PersistentFlags().GetBool("force-verify-sign"); err != nil {
		return nil, err
	}
	if noVerify, err = cmd.PersistentFlags().GetBool("no-verify"); err != nil {
		noVerify = false
	}

	var fileName string
	fileName, err = getFileName(cmd, args)
	if err != nil {
		return nil, err
	}
	tb, err := readTBLN(fileName)
	if err != nil {
		return nil, err
	}
	if forceSign && len(tb.Signs) == 0 {
		return nil, fmt.Errorf("signature verification failure")
	}
	if !noSign && (len(tb.Signs) != 0) {
		err = verifyAllSignature(tb)
		if err != nil {
			return nil, err
		}
		log.Println("Signature verification successful")
	} else {
		if !noVerify && !tb.Verify() {
			if len(tb.Hashes) == 0 {
				return nil, fmt.Errorf("no hash")
			}
			return nil, fmt.Errorf("verification failure")
		}
		log.Println("Verification successful")
	}
	return tb, nil
}
