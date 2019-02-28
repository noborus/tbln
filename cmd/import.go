package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	// MySQL driver
	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	_ "github.com/noborus/tbln/db/mysql"

	// PostgreSQL driver
	_ "github.com/noborus/tbln/db/postgres"
	// SQLlite3 driver
	_ "github.com/noborus/tbln/db/sqlite3"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "tbln import database table",
	Long:  `tbln import database table`,
	Run: func(cmd *cobra.Command, args []string) {
		dbImport(args)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func dbImport(args []string) error {
	fileName := args[0]
	if dbName == "" {
		return fmt.Errorf("should be db name")
	}
	conn, err := db.Open(dbName, dsn)
	if err != nil {
		log.Fatal(err)
	}
	conn.Tx, err = conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = db.WriteTable(conn, at, true)
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
