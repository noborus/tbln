package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [flags]",
	Short: "Merge two TBLNs",
	Long: `Merge two TBLNs.

	The two TBLNs should basically have the same structure.`,
	RunE: merge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.PersistentFlags().StringP("from-file", "", "", "TBLN FROM File")
	mergeCmd.PersistentFlags().StringP("to-file", "", "", "TBLN TO File")
	mergeCmd.PersistentFlags().StringP("from-dbname", "", "", "FROM database name")
	mergeCmd.PersistentFlags().StringP("from-dsn", "", "", "FROM dsn name")
	mergeCmd.PersistentFlags().StringP("from-schema", "", "", "FROM schema Name")
	mergeCmd.PersistentFlags().StringP("from-table", "", "", "FROM table name")
	mergeCmd.PersistentFlags().StringP("to-dbname", "", "", "TO database name")
	mergeCmd.PersistentFlags().StringP("to-dsn", "", "", "TO dsn name")
	mergeCmd.PersistentFlags().StringP("to-schema", "", "", "To schema Name")
	mergeCmd.PersistentFlags().StringP("to-table", "", "", "TO table name")
	mergeCmd.PersistentFlags().StringP("conflict", "", "no", `merge mode when conflict
 ignore - Ignores at insert conflict
 merge - insert update when conflicted
 sync - Synchronize (also execute delete)
 `)
	mergeCmd.PersistentFlags().BoolP("no-import", "", false, "output without importing into database")
	mergeCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")
}

func merge(cmd *cobra.Command, args []string) error {
	var err error
	var ffReader *tbln.Reader
	var fdReader, tdReader *db.Reader
	var fromFileName, toFileName string
	if len(args) >= 2 {
		fromFileName = args[0]
		toFileName = args[1]
	}
	if fromFileName, err = cmd.PersistentFlags().GetString("from-file"); err != nil {
		return err
	}
	if fromFileName != "" {
		dst, err := os.Open(fromFileName)
		if err != nil {
			return err
		}
		ffReader = tbln.NewReader(dst)
	}
	var fromDbName, fromDsn string
	if fromDbName, err = cmd.PersistentFlags().GetString("from-dbname"); err != nil {
		return err
	}
	if fromDbName != "" {
		if fromDsn, err = cmd.PersistentFlags().GetString("from-dsn"); err != nil {
			return err
		}
		fromConn, err := db.Open(fromDbName, fromDsn)
		if err != nil {
			return fmt.Errorf("%s: %s", fromDbName, err)
		}
		err = fromConn.Begin()
		if err != nil {
			return fmt.Errorf("%s: %s", fromDbName, err)
		}
		var schema string
		if schema, err = cmd.PersistentFlags().GetString("from-schema"); err != nil {
			return err
		}
		var tableName string
		if tableName, err = cmd.PersistentFlags().GetString("from-table"); err != nil {
			return err
		}
		fdReader, err = fromConn.ReadTable(schema, tableName, nil)
		if err != nil {
			return err
		}
	}
	from := decisionReader(ffReader, fdReader)
	var toDbName, toDsn string
	if toDbName, err = cmd.PersistentFlags().GetString("to-dbname"); err != nil {
		return err
	}
	if toDbName != "" {
		if toDsn, err = cmd.PersistentFlags().GetString("to-dsn"); err != nil {
			return err
		}
		toConn, err := db.Open(toDbName, toDsn)
		err = toConn.Begin()
		if err != nil {
			return fmt.Errorf("%s: %s", toDbName, err)
		}
		var schema string
		if schema, err = cmd.PersistentFlags().GetString("to-schema"); err != nil {
			return err
		}
		var tableName string
		if tableName, err = cmd.PersistentFlags().GetString("to-table"); err != nil {
			return err
		}
		var noImport bool
		if noImport, err = cmd.PersistentFlags().GetBool("no-import"); err != nil {
			return err
		}
		if !noImport {
			err = toConn.MergeTable(schema, tableName, from, false)
			if err != nil {
				log.Printf("Table Not found. Create Table [%s]\n", tableName)
				return err
			}
			return toConn.Commit()
		}
		tdReader, err = toConn.ReadTable(schema, tableName, nil)
		if err != nil {
			return err
		}
	}
	var tb *tbln.Tbln
	var tfReader *tbln.Reader
	if toFileName, err = cmd.PersistentFlags().GetString("to-file"); err != nil {
		return err
	}
	if toFileName != "" {
		src, err := os.Open(toFileName)
		if err != nil {
			return err
		}
		tfReader = tbln.NewReader(src)
	}
	to := decisionReader(tfReader, tdReader)
	tb, err = tbln.MergeAll(to, from)
	if err != nil {
		return err
	}
	return outputFile(tb, cmd)

}

func decisionReader(f *tbln.Reader, d *db.Reader) tbln.Comparer {
	if f != nil {
		return f
	}
	return d
}
