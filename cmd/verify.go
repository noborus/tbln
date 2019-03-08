package cmd

import (
	"fmt"
	"os"

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
	verifyCmd.PersistentFlags().BoolP("forcesign", "", false, "Force signature verification")
	verifyCmd.PersistentFlags().BoolP("nosign", "", false, "ignore signature verify")
	verifyCmd.PersistentFlags().StringP("pub", "p", "", "public Key File")
	verifyCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(verifyCmd)
}

func verify(cmd *cobra.Command, args []string) error {
	var err error
	var fileName, pubFileName string
	var pub []byte
	var forcesign, nosign bool
	if len(args) <= 0 {
		cmd.SilenceUsage = false
		return fmt.Errorf("require filename")
	}
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if fileName == "" && len(args) > 0 {
		fileName = args[0]
	}
	if nosign, err = cmd.PersistentFlags().GetBool("nosign"); err != nil {
		return err
	}
	if forcesign, err = cmd.PersistentFlags().GetBool("forcesign"); err != nil {
		return err
	}
	if pubFileName, err = cmd.PersistentFlags().GetString("pub"); err != nil {
		cmd.SilenceUsage = false
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
	if forcesign || (!nosign && (len(at.Signs) != 0)) {
		if len(pubFileName) == 0 {
			return fmt.Errorf("must be public key file")
		}
		pub, err = getPublicKey(pubFileName)
		if err != nil {
			return err
		}
		if at.VerifySignature(pub) {
			fmt.Println("Signature verified")
		} else {
			return fmt.Errorf("Signature verfication failure")
		}
	} else {
		if at.Verify() {
			fmt.Println("Verified")
		} else {
			return fmt.Errorf("Verfication failure")
		}
	}
	return nil
}
