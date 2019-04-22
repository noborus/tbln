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
	// SQLlite3 driver
	_ "github.com/noborus/tbln/db/sqlite3"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:          "import [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "Import database table",
	Long:         `Import the TBL file into the database.`,
	RunE:         dbImport,
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

func dbImport(cmd *cobra.Command, args []string) error {
	fileName, err := getFileName(cmd, args)
	if err != nil {
		return err
	}
	tb, err := verifiedTbln(cmd, args)
	if err != nil {
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
	conn, err := db.Open(u.Driver, u.DSN)
	if err != nil {
		return fmt.Errorf("%s: %s", u.Driver, err)
	}
	err = conn.Begin()
	if err != nil {
		return fmt.Errorf("%s: %s", u.Driver, err)
	}
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
	var cmode db.CreateMode
	if modeStr, err := cmd.PersistentFlags().GetString("mode"); err == nil {
		switch strings.ToLower(modeStr) {
		case "no":
			cmode = db.NotCreate
		case "create":
			cmode = db.Create
		case "ifnot":
			cmode = db.IfNotExists
		case "recreate":
			cmode = db.ReCreate
		case "only":
			cmode = db.CreateOnly
		default:
			cmode = db.Create
		}
	}
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
	var imode db.InsertMode
	switch {
	case ignore:
		imode = db.OrIgnore
	case update:
		imode = db.Update
	case delete:
		imode = db.Sync
	default:
		imode = db.Normal
	}
	err = writeImport(conn, tb, schema, cmode, imode)
	if err != nil {
		return err
	}
	err = conn.Tx.Commit()
	if err != nil {
		return err
	}
	log.Printf("[%s] import success\n", tb.TableName())
	return conn.Close()
}

func writeImport(conn *db.TDB, tb *tbln.Tbln, schema string, cmode db.CreateMode, imode db.InsertMode) error {
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

func mergeImport(conn *db.TDB, tb *tbln.Tbln, schema string, imode db.InsertMode) error {
	delete := false
	if imode == db.Sync {
		delete = true
	}
	return conn.MergeTableTbln(schema, tb.TableName(), tb, delete)
}
