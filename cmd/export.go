package cmd

import (
	"fmt"
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
	Use:          "export",
	SilenceUsage: true,
	Short:        "export database table or query",
	Long:         `export database table or query`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return dbExport(cmd, args)
	},
}

func init() {
	exportCmd.PersistentFlags().StringP("Schema", "n", "", "Schema Name")
	exportCmd.PersistentFlags().StringP("Table", "t", "", "Table Name")
	exportCmd.PersistentFlags().StringP("Query", "q", "", "SQL Query")
	exportCmd.PersistentFlags().StringP("sum", "c", "sha256", "CheckSum(sha256/sha512)")
	exportCmd.PersistentFlags().BoolP("sign", "s", false, "Sign TBLN file")
	exportCmd.PersistentFlags().StringP("key", "k", "", "Key File")

	rootCmd.AddCommand(exportCmd)
}

func dbExport(cmd *cobra.Command, args []string) error {
	var err error
	var schema string
	var tableName string
	var query string
	var sum string
	var keyFile string
	var priv []byte
	var signF bool
	if schema, err = cmd.PersistentFlags().GetString("Schema"); err != nil {
		return err
	}
	if tableName, err = cmd.PersistentFlags().GetString("Table"); err != nil {
		return err
	}
	if query, err = cmd.PersistentFlags().GetString("Query"); err != nil {
		return err
	}
	if signF, err = cmd.PersistentFlags().GetBool("sign"); err != nil {
		return err
	}
	if keyFile, err = cmd.PersistentFlags().GetString("key"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if keyFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be keyFile")
	}

	if tableName == "" && len(args) >= 1 {
		tableName = args[0]
	}
	if tableName == "" && query == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be Table Name or SQL Query")
	}
	if dbName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be Database Driver Name")
	}

	if sum, err = cmd.PersistentFlags().GetString("sum"); err != nil {
		return err
	}
	var sumHash tbln.HashType
	if sum == "sha512" {
		sumHash = tbln.SHA512
	} else {
		sumHash = tbln.SHA256
	}

	if signF {
		priv, err = decryptPrompt(keyFile)
		if err != nil {
			return err
		}
	}

	conn, err := db.Open(dbName, dsn)
	if err != nil {
		return err
	}
	var at *tbln.Tbln
	if query != "" {
		at, err = db.ReadQueryAll(conn, query)
	} else {
		at, err = db.ReadTableAll(conn, schema, tableName)
	}
	if err != nil {
		return err
	}
	err = at.SumHash(sumHash)
	if err != nil {
		return err
	}
	if signF {
		at.Sign(priv)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}
	return conn.Close()
}
