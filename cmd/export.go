package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
	srcdbName string
	srcdsn    string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:          "export  [flags] <Table Name> or <SQL Query>",
	SilenceUsage: true,
	Short:        "export database table or query",
	Long: `Export from the database by table or SQL Query.
Add as much of the table information as possible including 
the column name and column type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return dbExport(cmd, args)
	},
}

func init() {
	exportCmd.PersistentFlags().StringVar(&srcdbName, "db", "", "database name")
	exportCmd.PersistentFlags().StringVar(&srcdsn, "dsn", "", "dsn name")
	exportCmd.PersistentFlags().StringP("Schema", "n", "", "Schema Name")
	exportCmd.PersistentFlags().StringP("Table", "t", "", "Table Name")
	exportCmd.PersistentFlags().StringP("SQL", "", "", "SQL Query")
	exportCmd.PersistentFlags().StringP("hash", "a", "sha256", "Hash algorithm(sha256 or sha512)")
	exportCmd.PersistentFlags().BoolP("sign", "", false, "Sign TBLN file")
	exportCmd.PersistentFlags().StringP("key", "k", "", "Key File")

	rootCmd.AddCommand(exportCmd)
}

func dbExport(cmd *cobra.Command, args []string) error {
	var err error
	var schema string
	var tableName string
	var sql string
	var sumHash string
	var signF bool
	var keyFile string
	if schema, err = cmd.PersistentFlags().GetString("Schema"); err != nil {
		return err
	}
	if tableName, err = cmd.PersistentFlags().GetString("Table"); err != nil {
		return err
	}
	if sql, err = cmd.PersistentFlags().GetString("SQL"); err != nil {
		return err
	}
	if signF, err = cmd.PersistentFlags().GetBool("sign"); err != nil {
		return err
	}
	if keyFile, err = cmd.PersistentFlags().GetString("key"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if signF && keyFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be keyFile")
	}

	if tableName == "" && len(args) >= 1 {
		tableName = args[0]
	}
	if tableName == "" && sql == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be Table Name or SQL Query")
	}
	if srcdbName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be Database Driver Name")
	}

	if sumHash, err = cmd.PersistentFlags().GetString("hash"); err != nil {
		return err
	}

	var privKey []byte
	if signF {
		privKey, err = decryptPrompt(keyFile)
		if err != nil {
			return err
		}
	}
	keyName := filepath.Base(keyFile[:len(keyFile)-len(filepath.Ext(keyFile))])

	conn, err := db.Open(srcdbName, srcdsn)
	if err != nil {
		return fmt.Errorf("%s: %s", srcdbName, err)
	}
	var at *tbln.Tbln
	if sql != "" {
		at, err = db.ReadQueryAll(conn, sql)
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
		at.Sign(keyName, privKey)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}
	return conn.Close()
}
