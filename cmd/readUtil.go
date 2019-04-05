package cmd

import (
	"fmt"
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

func readTbln(fName string, cmd *cobra.Command) (*tbln.Tbln, error) {
	var err error
	var fileName string
	if fileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
		cmd.SilenceUsage = false
		return nil, err
	}
	if fileName == "" {
		fileName = fName
	}
	if fileName == "" {
		cmd.SilenceUsage = false
		return nil, fmt.Errorf("require filename")
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
	return at, nil
}
