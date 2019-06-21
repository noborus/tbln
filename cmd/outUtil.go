package cmd

import (
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

func outputFile(at *tbln.TBLN, cmd *cobra.Command) error {
	var err error
	var output string
	if output, err = cmd.PersistentFlags().GetString("output"); err != nil {
		return err
	}
	var out *os.File
	if output == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(output)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return tbln.WriteAll(out, at)
}
