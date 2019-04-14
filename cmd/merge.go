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
	mergeCmd.PersistentFlags().StringP("from-db", "", "", "FROM database driver name")
	mergeCmd.PersistentFlags().StringP("from-dsn", "", "", "FROM dsn name")
	mergeCmd.PersistentFlags().StringP("from-schema", "", "", "FROM schema Name")
	mergeCmd.PersistentFlags().StringP("from-table", "", "", "FROM table name")
	mergeCmd.PersistentFlags().StringP("to-db", "", "", "TO database driver name")
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
		var conflict string
		if conflict, err = cmd.PersistentFlags().GetString("conflict"); err != nil {
			return err
		}
		var imode db.InsertMode
		switch conflict {
		case "ignore":
			imode = db.OrIgnore
		case "merge":
			imode = db.Merge
		case "sync":
			imode = db.Sync
		default:
			imode = db.Normal
		}

		var noImport bool
		if noImport, err = cmd.PersistentFlags().GetBool("no-import"); err != nil {
			return err
		}
		if !noImport {
			if imode != db.OrIgnore {
				delete := false
				if imode == db.Sync {
					delete = true
				}
				err = toConn.MergeTable(schema, tableName, fromReader, delete)
				if err == nil {
					return toConn.Commit()
				}
				log.Printf("Table Not found. Create Table [%s]\n", tableName)
			}
			err = toConn.WriteReader(fromReader, schema, tableName, db.Create, db.OrIgnore)
			if err != nil {
				return err
			}
			return toConn.Commit()
		}
		toReader, err = toConn.ReadTable(schema, tableName, nil)
		if err != nil {
			return err
		}
	}
	var tb *tbln.Tbln
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
	tb, err = tbln.MergeAll(toReader, fromReader)
	if err != nil {
		return err
	}
	return outputFile(tb, cmd)

}
