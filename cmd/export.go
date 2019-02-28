package cmd

import (
	"fmt"
	"log"
	"os"

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

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "tbln export database table",
	Long:  `tbln export database table`,
	Run: func(cmd *cobra.Command, args []string) {
		dbExport(args)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func dbExport(args []string) error {
	tableName := args[0]
	if dbName == "" {
		return fmt.Errorf("should be db name")
	}
	conn, err := db.Open(dbName, dsn)
	if err != nil {
		log.Fatal(err)
	}
	at, err := db.ReadTableAll(conn, tableName)
	if err != nil {
		log.Fatal(err)
	}
	comment := fmt.Sprintf("DB:%s\tTable:%s", conn.Name, tableName)
	at.Comments = []string{comment}
	_, err = at.SumHash(tbln.SHA256)
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
	return conn.Close()
}
