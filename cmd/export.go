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
	exportCmd.PersistentFlags().StringP("schema", "n", "", "Schema Name")
	exportCmd.PersistentFlags().StringP("table", "t", "", "Table Name")
	exportCmd.PersistentFlags().StringP("query", "", "", "SQL Query")
	exportCmd.PersistentFlags().StringP("hash", "a", "sha256", "Hash algorithm(sha256 or sha512)")
	exportCmd.PersistentFlags().BoolP("sign", "", false, "Sign TBLN file")
	exportCmd.PersistentFlags().BoolP("quiet", "q", false, "Do not prompt for password.")

	rootCmd.AddCommand(exportCmd)
}

func dbExport(cmd *cobra.Command, args []string) error {
	var err error

	var schema string
	var tableName string
	var query string
	var sumHash string

	var signF bool

	if schema, err = cmd.PersistentFlags().GetString("schema"); err != nil {
		return err
	}
	if tableName, err = cmd.PersistentFlags().GetString("table"); err != nil {
		return err
	}
	if query, err = cmd.PersistentFlags().GetString("query"); err != nil {
		return err
	}
	if signF, err = cmd.PersistentFlags().GetBool("sign"); err != nil {
		return err
	}
	if signF && seckey == "" {
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
	if srcdbName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("must be Database Driver Name")
	}
	var quiet bool
	if quiet, err = cmd.PersistentFlags().GetBool("quiet"); err != nil {
		cmd.SilenceUsage = false
		return err
	}
	if sumHash, err = cmd.PersistentFlags().GetString("hash"); err != nil {
		return err
	}

	conn, err := db.Open(srcdbName, srcdsn)
	if err != nil {
		return fmt.Errorf("%s: %s", srcdbName, err)
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
		privKey, err := getPrivateKeyFile(seckey, keyname)
		if err != nil {
			return err
		}
		var priv []byte
		if quiet {
			priv, err = decrypt([]byte(""), privKey)
		} else {
			priv, err = decryptPrompt(privKey)
		}
		if err != nil {
			return err
		}
		at.Sign(keyname, priv)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		return err
	}
	return conn.Close()
}
