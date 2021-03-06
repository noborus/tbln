package cmd

import (
	"fmt"
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:          "diff [flags] [file or DB] [file or DB]",
	SilenceUsage: true,
	Short:        "Diff two TBLNs",
	Long: `Diff two TBLNs

The two TBLNs should basically have the same structure.`,
	RunE: diff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.PersistentFlags().StringP("file", "f", "", "self TBLN file")
	diffCmd.PersistentFlags().StringP("dburl", "d", "", "self database url")
	diffCmd.PersistentFlags().StringP("schema", "n", "", "self schema Name")
	diffCmd.PersistentFlags().StringP("table", "t", "", "self table name")
	diffCmd.PersistentFlags().StringP("other-file", "", "", "other TBLN file")
	diffCmd.PersistentFlags().StringP("other-dburl", "", "", "other database url")
	diffCmd.PersistentFlags().StringP("other-schema", "", "", "other schema Name")
	diffCmd.PersistentFlags().StringP("other-table", "", "", "other table name")
	diffCmd.PersistentFlags().BoolP("all", "", true, "show all, including unchanged rows")
	diffCmd.PersistentFlags().BoolP("change", "", false, "show only changed rows")
	diffCmd.PersistentFlags().BoolP("update", "", false, "show only rows to add and update")
	diffCmd.PersistentFlags().SortFlags = false
	diffCmd.Flags().SortFlags = false
}

func diff(cmd *cobra.Command, args []string) error {
	otherReader, err := getOtherReader(cmd, args)
	if err != nil {
		return err
	}
	selfReader, err := getSelfReader(cmd, args)
	if err != nil {
		return err
	}
	if otherReader == nil || selfReader == nil {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires other and self")
	}
	var change, update bool
	if change, err = cmd.PersistentFlags().GetBool("change"); err != nil {
		return err
	}
	if update, err = cmd.PersistentFlags().GetBool("update"); err != nil {
		return err
	}
	var flag tbln.DiffMode
	switch {
	case change:
		flag = tbln.OnlyDiff
	case update:
		flag = tbln.OnlyAdd
	default:
		flag = tbln.AllDiff
	}
	err = tbln.DiffAll(os.Stdout, otherReader, selfReader, flag)
	if err != nil {
		return err
	}
	return nil
}
