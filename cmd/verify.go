package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify tbln file",
	Long:  `verify tbln file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return verify(cmd, args)
	},
}

func init() {
	verifyCmd.PersistentFlags().StringP("pub", "p", "", "public Key File")
	verifyCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")

	rootCmd.AddCommand(verifyCmd)
}

func verify(cmd *cobra.Command, args []string) error {
	var err error
	var fileName, pubFileName string
	var pub []byte
	if len(args) <= 0 {
		return fmt.Errorf("require filename")
	}
	fileName = args[0]
	/*
		if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
			log.Fatal(err)
		}
	*/
	if pubFileName, err = cmd.PersistentFlags().GetString("pub"); err != nil {
		log.Fatal(err)
	}
	if len(pubFileName) > 0 {
		pubFile, _ := os.Open(pubFileName)
		defer pubFile.Close()
		fi, _ := pubFile.Stat()
		size := fi.Size()
		data := make([]byte, size)
		pubFile.Read(data)
		pub, _ = base64.StdEncoding.DecodeString(string(data))
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if len(pubFileName) > 0 {
		if at.VerifySignature(pub) {
			fmt.Println("Signature verified")
		} else {
			fmt.Println("Signature verfication failure")
			os.Exit(1)
		}
	} else {
		if at.Verify() {
			fmt.Println("Verified")
		} else {
			fmt.Println("Verfication failure")
			os.Exit(1)
		}
	}
	return nil
}
