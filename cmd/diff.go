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
	diffCmd.PersistentFlags().BoolP("all", "", false, "Show all, including unchanged rows")
	diffCmd.PersistentFlags().BoolP("change", "", false, "Show only changed rows")
	diffCmd.PersistentFlags().BoolP("update", "", false, "Show only rows to add and update")
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
	var all, change, update bool
	if all, err = cmd.PersistentFlags().GetBool("all"); err != nil {
		return err
	}
	if change, err = cmd.PersistentFlags().GetBool("change"); err != nil {
		return err
	}
	if update, err = cmd.PersistentFlags().GetBool("update"); err != nil {
		return err
	}
	var flag tbln.DiffMode
	switch {
	case all:
		flag = tbln.AllDiff
	case change:
		flag = tbln.OnlyDiff
	case update:
		flag = tbln.OnlyAdd
	default:
		flag = tbln.AllDiff
	}
	err = tbln.DiffAll(os.Stdout, fromReader, toReader, flag)
	if err != nil {
		return err
	}
	return nil
}
