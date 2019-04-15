package cmd

import (
	"fmt"
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:          "diff [flags]",
	SilenceUsage: true,
	Short:        "Diff two TBLNs",
	Long: `Diff two TBLNs

	The two TBLNs should basically have the same structure.`,
	RunE: diff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.PersistentFlags().StringP("from-file", "", "", "TBLN FROM File")
	diffCmd.PersistentFlags().StringP("to-file", "", "", "TBLN TO File")
	diffCmd.PersistentFlags().StringP("from-db", "", "", "FROM database driver name")
	diffCmd.PersistentFlags().StringP("from-dsn", "", "", "FROM dsn name")
	diffCmd.PersistentFlags().StringP("from-schema", "", "", "FROM schema Name")
	diffCmd.PersistentFlags().StringP("from-table", "", "", "FROM table name")
	diffCmd.PersistentFlags().StringP("to-db", "", "", "TO database driver name")
	diffCmd.PersistentFlags().StringP("to-dsn", "", "", "TO dsn name")
	diffCmd.PersistentFlags().StringP("to-schema", "", "", "To schema Name")
	diffCmd.PersistentFlags().StringP("to-table", "", "", "TO table name")
}

func diff(cmd *cobra.Command, args []string) error {
	fromReader, err := getFromReader(cmd, args)
	if err != nil {
		return err
	}
	toReader, err := getToReader(cmd, args)
	if err != nil {
		return err
	}
	if fromReader == nil || toReader == nil {
		return fmt.Errorf("requires from and to")
	}
	err = tbln.DiffAll(os.Stdout, toReader, fromReader)
	if err != nil {
		return err
	}
	return nil
}
