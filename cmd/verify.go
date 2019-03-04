package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		verify(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func verify(cmd *cobra.Command, args []string) error {
	fileName := args[0]
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
	if at.Verify() {
		fmt.Println("verify")
	} else {
		fmt.Println("verify error")
		os.Exit(1)
	}
	return nil
}
