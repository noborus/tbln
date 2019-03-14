package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/noborus/tbln/db"

	// MySQL driver
	_ "github.com/noborus/tbln/db/mysql"
	// PostgreSQL driver
	_ "github.com/noborus/tbln/db/postgres"
	// SQLlite3 driver
	_ "github.com/noborus/tbln/db/sqlite3"
)

// global variable from global flags
var (
	destdbName string
	destdsn    string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:          "import [flags] [TBLN file]",
	SilenceUsage: true,
	Short:        "import database table",
	Long: `Import the TBL file into the database.
<mode> specifies how to execute when there is already a table.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return dbImport(cmd, args)
	},
}

func init() {
	importCmd.PersistentFlags().StringVar(&destdbName, "db", "", "database name")
	importCmd.PersistentFlags().StringVar(&destdsn, "dsn", "", "dsn name")
	importCmd.PersistentFlags().StringP("Schema", "n", "", "Schema Name")
	importCmd.PersistentFlags().StringP("mode", "m", "i", "create mode (a:NotCreate/c:Create/i:IfNotExists/r:ReCreate/s:CreateOnly")
	importCmd.PersistentFlags().StringP("table", "t", "", "Table Name")
	importCmd.PersistentFlags().BoolP("force-verify-sign", "", false, "Force signature verification")
	importCmd.PersistentFlags().BoolP("no-verify-sign", "", false, "ignore signature verify")
	importCmd.PersistentFlags().StringP("file", "f", "", "TBLN File")
	rootCmd.AddCommand(importCmd)
}

func dbImport(cmd *cobra.Command, args []string) error {
	at, err := verifiedTbln(cmd, args)
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
	conn.Tx, err = conn.Begin()
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
		at.SetTableName(tableName)
	}
	var mode db.CreateMode
	if modeStr, err := cmd.PersistentFlags().GetString("mode"); err == nil {
		switch strings.ToLower(modeStr) {
		case "a":
			mode = db.NotCreate
		case "c":
			mode = db.Create
		case "i":
			mode = db.IfNotExists
		case "r":
			mode = db.ReCreate
		case "s":
			mode = db.CreateOnly
		default:
			mode = db.Create
		}
	}
	err = db.WriteTable(conn, at, schema, mode)
	if err != nil {
		return err
	}
	err = conn.Tx.Commit()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s import success\n", at.TableName())
	return conn.Close()
}
