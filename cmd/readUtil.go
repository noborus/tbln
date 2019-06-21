package cmd

import (
	"fmt"
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

func getFileName(cmd *cobra.Command, args []string) (string, error) {
	var err error
	var fileName string
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		cmd.SilenceUsage = false
		return "", err
	}
	if fileName == "" && len(args) > 0 {
		fileName = args[0]
	}
	if fileName == "" {
		cmd.SilenceUsage = false
		return "", fmt.Errorf("require filename")
	}
	return fileName, nil
}

func readTBLN(fileName string) (*tbln.TBLN, error) {
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
	return at, nil
}
