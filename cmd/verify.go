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
	_, err := verifiedTbln(cmd, args)
	return err
}

func verifiedTbln(cmd *cobra.Command, args []string) (*tbln.Tbln, error) {
	var err error
	var fileName string
	if fileName == "" && len(args) > 0 {
		fileName = args[0]
	}

	var forcesign, nosign, noverify bool
	if nosign, err = cmd.PersistentFlags().GetBool("no-verify-sign"); err != nil {
		return nil, err
	}
	if forcesign, err = cmd.PersistentFlags().GetBool("force-verify-sign"); err != nil {
		return nil, err
	}
	if noverify, err = cmd.PersistentFlags().GetBool("no-verify"); err != nil {
		noverify = false
	}

	at, err := readTbln(fileName, cmd)
	if err != nil {
		return nil, err
	}
	if forcesign && len(at.Signs) == 0 {
		return nil, fmt.Errorf("signature verification failure")
	}
	if !nosign && (len(at.Signs) != 0) {
		for kname := range at.Signs {
			pub, err := keystore.Search(KeyStore, kname)
			if err != nil {
				return nil, err
			}
			if !at.VerifySignature(kname, pub) {
				return nil, fmt.Errorf("signature verification failure")
			}
		}
		log.Println("Signature verification successful")
	} else {
		if !noverify && !at.Verify() {
			if len(at.Hashes) == 0 {
				return nil, fmt.Errorf("no hash")
			}
			return nil, fmt.Errorf("verification failure")
		}
		log.Println("Verification successful")
	}
	return at, nil
}
