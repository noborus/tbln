package cmd

import (
	"fmt"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// exceptCmd represents the except command
var exceptCmd = &cobra.Command{
	Use:          "except",
	SilenceUsage: true,
	Short:        "Difference set two TBLNs",
	Long: `Difference set two TBLNs

	The two TBLNs should basically have the same structure.`,
	RunE: except,
}

func init() {
	rootCmd.AddCommand(exceptCmd)
	exceptCmd.PersistentFlags().StringP("from-file", "", "", "TBLN FROM File")
	exceptCmd.PersistentFlags().StringP("to-file", "", "", "TBLN TO File")
	exceptCmd.PersistentFlags().StringP("from-db", "", "", "FROM database driver name")
	exceptCmd.PersistentFlags().StringP("from-dsn", "", "", "FROM dsn name")
	exceptCmd.PersistentFlags().StringP("from-schema", "", "", "FROM schema Name")
	exceptCmd.PersistentFlags().StringP("from-table", "", "", "FROM table name")
	exceptCmd.PersistentFlags().StringP("to-db", "", "", "TO database driver name")
	exceptCmd.PersistentFlags().StringP("to-dsn", "", "", "TO dsn name")
	exceptCmd.PersistentFlags().StringP("to-schema", "", "", "To schema Name")
	exceptCmd.PersistentFlags().StringP("to-table", "", "", "TO table name")
	exceptCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")
}

func except(cmd *cobra.Command, args []string) error {
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
	var tb *tbln.Tbln
	tb, err = tbln.ExceptAll(toReader, fromReader)
	if err != nil {
		return err
	}
	return outputFile(tb, cmd)
}
