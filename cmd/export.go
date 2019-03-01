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
	Short: "export database table",
	Long:  `export database table`,
	Run: func(cmd *cobra.Command, args []string) {
		dbExport(cmd, args)
	},
}

func init() {
	exportCmd.PersistentFlags().StringP("Table", "t", "", "Table Name")
	exportCmd.PersistentFlags().StringP("Query", "q", "", "SQL Query")
	exportCmd.PersistentFlags().StringP("sum", "s", "sha256", "CheckSum(sha256/sha512)")
	rootCmd.AddCommand(exportCmd)
}

func dbExport(cmd *cobra.Command, args []string) error {
	var err error
	var tableName string
	var query string
	if tableName, err = cmd.PersistentFlags().GetString("Table"); err != nil {
		log.Fatal(err)
	}
	if query, err = cmd.PersistentFlags().GetString("Query"); err != nil {
		log.Fatal(err)
	}
	sum := "sha256"
	if sum, err = cmd.PersistentFlags().GetString("sum"); err != nil {
		log.Fatal(err)
	}
	var sumHash tbln.HashType
	if sum == "sha512" {
		sumHash = tbln.SHA512
	} else {
		sumHash = tbln.SHA256
	}
	if tableName == "" && len(args) >= 1 {
		tableName = args[0]
	}
	if tableName == "" && query == "" {
		log.Fatal("should be table name or query")
	}
	if dbName == "" {
		log.Fatal("should be db name")
	}

	conn, err := db.Open(dbName, dsn)
	if err != nil {
		log.Fatal(err)
	}
	var at *tbln.Tbln
	if query != "" {
		at, err = db.ReadQueryAll(conn, query)
	} else {
		at, err = db.ReadTableAll(conn, tableName)
	}
	if err != nil {
		log.Fatal(err)
	}
	comment := fmt.Sprintf("DB:%s\tTable:%s", conn.Name, tableName)
	at.Comments = []string{comment}
	_, err = at.SumHash(sumHash)
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
	return conn.Close()
}
