package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"

	"github.com/spf13/cobra"

	// MySQL driver
	_ "github.com/noborus/tbln/db/mysql"
	// PostgreSQL driver
	_ "github.com/noborus/tbln/db/postgres"
	// SQLlite3 driver
	_ "github.com/noborus/tbln/db/sqlite3"
)

// database variable
var (
	destdbName string
	destdsn    string
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
	importCmd.PersistentFlags().StringVar(&destdbName, "db", "", "database name")
	importCmd.PersistentFlags().StringVar(&destdsn, "dsn", "", "dsn name")
	importCmd.PersistentFlags().StringP("Schema", "n", "", "schema Name")
	importCmd.PersistentFlags().StringP("mode", "m", "create", `create mode
 no		- Insert without creating a table.
 create	- Normal create.
 ifnot	- If the table does not exist, it will be created.
 recreate	- Drop and re-create the table.
 only	- Only create a table, do not insert data.
 `)
	importCmd.PersistentFlags().StringP("conflict", "", "no", `merge mode when conflict
 ignore - Ignores at insert conflict
 merge - insert update when conflicted
 sync - Synchronize (also execute delete)
 `)
	importCmd.PersistentFlags().StringP("table", "t", "", "table name")
	importCmd.PersistentFlags().BoolP("force-verify-sign", "", false, "force signature verification")
	importCmd.PersistentFlags().BoolP("no-verify-sign", "", false, "ignore signature verify")
	importCmd.PersistentFlags().BoolP("no-verify", "", false, "ignore verify")
	importCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(importCmd)
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
	if destdbName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires database driver name")
	}
	conn, err := db.Open(destdbName, destdsn)
	if err != nil {
		return fmt.Errorf("%s: %s", destdbName, err)
	}
	err = conn.Begin()
	if err != nil {
		return fmt.Errorf("%s: %s", destdbName, err)
	}
	var schema string
	if schema, err = cmd.PersistentFlags().GetString("Schema"); err != nil {
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
	if imode >= db.Merge {
		err := mergeImport(conn, tb, schema, imode)
		if err == nil {
			return nil
		}
		log.Printf("Table Not found. Create Table [%s]\n", tb.TableName())
	}
	return db.WriteTable(conn, tb, schema, cmode, imode)
}

func mergeImport(conn *db.TDB, tb *tbln.Tbln, schema string, imode db.InsertMode) error {
	sr := tbln.NewSelfReader(tb)
	delete := false
	if imode == db.Sync {
		delete = true
	}
	return conn.MergeTable(schema, tb.TableName(), sr, delete)
}
