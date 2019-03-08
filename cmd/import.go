package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/noborus/tbln"
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
	Long:         `import database table`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return dbImport(cmd, args)
	},
}

func init() {
	importCmd.PersistentFlags().StringVar(&destdbName, "db", "", "database name")
	importCmd.PersistentFlags().StringVar(&destdsn, "dsn", "", "dsn name")
	importCmd.PersistentFlags().StringP("mode", "m", "i", "create mode (a:NotCreate/c:Create/i:IfNotExists/r:ReCreate/s:CreateOnly")
	importCmd.PersistentFlags().StringP("table", "t", "", "Table Name")
	rootCmd.AddCommand(importCmd)
}

func dbImport(cmd *cobra.Command, args []string) error {
	var fileName string
	if len(args) > 0 {
		fileName = args[0]
	}
	if destdbName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be database name")
	}
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("%s: %s", fileName, err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%s: %s", fileName, err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("%s: %s", fileName, err)
	}

	conn, err := db.Open(destdbName, destdsn)
	if err != nil {
		return fmt.Errorf("%s: %s", destdbName, err)
	}
	conn.Tx, err = conn.Begin()
	if err != nil {
		return fmt.Errorf("%s: %s", destdbName, err)
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
	err = db.WriteTable(conn, at, mode)
	if err != nil {
		return err
	}
	err = conn.Tx.Commit()
	if err != nil {
		return err
	}
	fmt.Printf("%s import success\n", fileName)
	return conn.Close()
}
