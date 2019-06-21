package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	"github.com/xo/dburl"

	"github.com/spf13/cobra"

	// MySQL driver
	_ "github.com/noborus/tbln/db/mysql"
	// PostgreSQL driver
	_ "github.com/noborus/tbln/db/postgres"
	// SQLite3 driver
	_ "github.com/noborus/tbln/db/sqlite3"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:          "import [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Import database table",
	Long:         `Import the TBL file into the database.`,
	RunE:         cmdImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().StringP("dburl", "d", "", "database url")
	importCmd.PersistentFlags().StringP("schema", "n", "", "schema Name")
	importCmd.PersistentFlags().StringP("mode", "m", "ifnot", `create mode
 no		- Insert without creating a table.
 create	- Normal create.
 ifnot	- If the table does not exist, it will be created.
 recreate	- Drop and re-create the table.
 only	- Only create a table, do not insert data.
 `)
	importCmd.PersistentFlags().BoolP("ignore", "", false, "Ignore when conflict occurs")
	importCmd.PersistentFlags().BoolP("update", "", false, "Update when a conflict occurs")
	importCmd.PersistentFlags().BoolP("delete", "", false, "Execute update and delete rows not in other")
	importCmd.PersistentFlags().StringP("table", "t", "", "table name")
	importCmd.PersistentFlags().BoolP("force-verify-sign", "", false, "force signature verification")
	importCmd.PersistentFlags().BoolP("no-verify-sign", "", false, "ignore signature verify")
	importCmd.PersistentFlags().BoolP("no-verify", "", false, "ignore verify")
	importCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	importCmd.PersistentFlags().SortFlags = false
	importCmd.Flags().SortFlags = false
}

func createMode(mode string) db.CreateMode {
	switch strings.ToLower(mode) {
	case "no":
		return db.NotCreate
	case "create":
		return db.Create
	case "ifnot":
		return db.IfNotExists
	case "recreate":
		return db.ReCreate
	case "only":
		return db.CreateOnly
	default:
		return db.Create
	}
}

func insertMode(ignore bool, update bool, delete bool) db.InsertMode {
	switch {
	case ignore:
		return db.OrIgnore
	case update:
		return db.Update
	case delete:
		return db.Sync
	default:
		return db.Normal
	}
}

func cmdImport(cmd *cobra.Command, args []string) error {
	fileName, err := getFileName(cmd, args)
	if err != nil {
		return err
	}
	tb, err := verifiedTBLN(cmd, args)
	if err != nil {
		log.Printf("invalid tbln file [%s]", fileName)
		return err
	}
	var url string
	if url, err = cmd.PersistentFlags().GetString("dburl"); err != nil {
		return err
	}
	if url == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("database url not found")
	}
	u, err := dburl.Parse(url)
	if err != nil {
		return fmt.Errorf("%s: %s", url, err)
	}

	var modeStr string
	if modeStr, err = cmd.PersistentFlags().GetString("mode"); err != nil {
		return err
	}
	cmode := createMode(modeStr)

	var ignore, update, delete bool
	if ignore, err = cmd.PersistentFlags().GetBool("ignore"); err != nil {
		return err
	}
	if update, err = cmd.PersistentFlags().GetBool("update"); err != nil {
		return err
	}
	if delete, err = cmd.PersistentFlags().GetBool("delete"); err != nil {
		return err
	}
	imode := insertMode(ignore, update, delete)

	var schema string
	if schema, err = cmd.PersistentFlags().GetString("schema"); err != nil {
		return err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("table"); err != nil {
		return err
	}
	if tableName != "" {
		tb.SetTableName(tableName)
	}
	if tb.TableName() == "" {
		base := filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])
		tb.SetTableName(base)
	}
	return dbImport(u, tb, schema, cmode, imode)
}

func dbImport(u *dburl.URL, tb *tbln.TBLN, schema string, cmode db.CreateMode, imode db.InsertMode) error {
	conn, err := db.Open(u.Driver, u.DSN)
	if err != nil {
		return fmt.Errorf("%s: %s", u.Driver, err)
	}
	err = conn.Begin()
	if err != nil {
		return fmt.Errorf("%s: %s", u.Driver, err)
	}
	err = writeImport(conn, tb, schema, cmode, imode)
	if err != nil {
		log.Printf("[%s] import failure\n", tb.TableName())
		return err
	}
	err = conn.Tx.Commit()
	if err != nil {
		log.Printf("[%s] import failure\n", tb.TableName())
		return err
	}
	log.Printf("[%s] import success\n", tb.TableName())
	return conn.Close()
}

func writeImport(conn *db.TDB, tb *tbln.TBLN, schema string, cmode db.CreateMode, imode db.InsertMode) error {
	if imode >= db.Update {
		err := mergeImport(conn, tb, schema, imode)
		if err == nil {
			return nil
		}
		log.Printf("ERR:%s", err)
		log.Printf("Table Not found. Create Table [%s]\n", tb.TableName())
	}
	return db.WriteTable(conn, tb, schema, cmode, imode)
}

func mergeImport(conn *db.TDB, tb *tbln.TBLN, schema string, imode db.InsertMode) error {
	delete := false
	if imode == db.Sync {
		delete = true
	}
	return conn.MergeTableTBLN(schema, tb.TableName(), tb, delete)
}
