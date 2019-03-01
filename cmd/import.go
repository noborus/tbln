package cmd

import (
	"fmt"
	"log"
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

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import database table",
	Long:  `import database table`,
	Run: func(cmd *cobra.Command, args []string) {
		dbImport(cmd, args)
	},
}

func init() {
	importCmd.PersistentFlags().StringP("mode", "m", "i", "create mode (n:NotCreate/c:Create/i:IfNotExists/r:ReCreate")
	importCmd.PersistentFlags().StringP("table", "t", "", "Table Name")
	rootCmd.AddCommand(importCmd)
}

func dbImport(cmd *cobra.Command, args []string) error {
	fileName := args[0]
	if dbName == "" {
		return fmt.Errorf("should be db name")
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.Open(dbName, dsn)
	if err != nil {
		log.Fatal(err)
	}
	conn.Tx, err = conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	var mode db.CreateMode
	if modeStr, err := cmd.PersistentFlags().GetString("mode"); err == nil {
		switch strings.ToLower(modeStr) {
		case "n":
			mode = db.NotCreate
		case "c":
			mode = db.Create
		case "i":
			mode = db.IfNotExists
		case "r":
			mode = db.ReCreate
		default:
			mode = db.Create
		}
	}
	err = db.WriteTable(conn, at, mode)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("import called")
	return conn.Close()
}
