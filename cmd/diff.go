package cmd

import (
	"fmt"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
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
	var err error
	var fromReader, toReader tbln.Reader

	var fromFileName, toFileName string
	if fromFileName, err = cmd.PersistentFlags().GetString("from-file"); err != nil {
		return err
	}
	if fromFileName != "" {
		dst, err := os.Open(fromFileName)
		if err != nil {
			return err
		}
		fromReader = tbln.NewReader(dst)
	}
	var fromDb, fromDsn string
	if fromDb, err = cmd.PersistentFlags().GetString("from-db"); err != nil {
		return err
	}
	if fromDb != "" {
		if fromDsn, err = cmd.PersistentFlags().GetString("from-dsn"); err != nil {
			return err
		}
		fromConn, err := db.Open(fromDb, fromDsn)
		if err != nil {
			return fmt.Errorf("%s: %s", fromDb, err)
		}
		err = fromConn.Begin()
		if err != nil {
			return fmt.Errorf("%s: %s", fromDb, err)
		}
		var schema string
		if schema, err = cmd.PersistentFlags().GetString("from-schema"); err != nil {
			return err
		}
		var tableName string
		if tableName, err = cmd.PersistentFlags().GetString("from-table"); err != nil {
			return err
		}
		fromReader, err = fromConn.ReadTable(schema, tableName, nil)
		if err != nil {
			return err
		}
	}

	if toFileName, err = cmd.PersistentFlags().GetString("to-file"); err != nil {
		return err
	}
	if toFileName != "" {
		src, err := os.Open(toFileName)
		if err != nil {
			return err
		}
		toReader = tbln.NewReader(src)
	}
	var toDb, toDsn string
	if toDb, err = cmd.PersistentFlags().GetString("to-db"); err != nil {
		return err
	}
	if toDb != "" {
		if toDsn, err = cmd.PersistentFlags().GetString("to-dsn"); err != nil {
			return err
		}
		toConn, err := db.Open(toDb, toDsn)
		if err != nil {
			return fmt.Errorf("%s: %s", toDb, err)
		}
		err = toConn.Begin()
		if err != nil {
			return fmt.Errorf("%s: %s", toDb, err)
		}
		var schema string
		if schema, err = cmd.PersistentFlags().GetString("to-schema"); err != nil {
			return err
		}
		var tableName string
		if tableName, err = cmd.PersistentFlags().GetString("to-table"); err != nil {
			return err
		}
		toReader, err = toConn.ReadTable(schema, tableName, nil)
		if err != nil {
			return err
		}
	}
	err = tbln.DiffAll(os.Stdout, toReader, fromReader)
	if err != nil {
		return err
	}
	return nil
}
