package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/cmd/key"
	"github.com/noborus/tbln/db"

	// MySQL driver
	_ "github.com/noborus/tbln/db/mysql"
	// PostgreSQL driver
	_ "github.com/noborus/tbln/db/postgres"
	// SQLlite3 driver
	_ "github.com/noborus/tbln/db/sqlite3"
)

// database variable
var (
	srcdbName string
	srcdsn    string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:          "export [flags] <Table Name> or <SQL Query>",
	SilenceUsage: true,
	Short:        "Export database table or query",
	Long: `Export from the database by table or SQL Query.
Add as much of the table information as possible including 
the column name and column type.`,
	RunE: dbExport,
}

func init() {
	exportCmd.PersistentFlags().StringVar(&srcdbName, "db", "", "database name")
	exportCmd.PersistentFlags().StringVar(&srcdsn, "dsn", "", "dsn name")
	exportCmd.PersistentFlags().StringP("schema", "n", "", "schema name")
	exportCmd.PersistentFlags().StringP("table", "t", "", "table name")
	exportCmd.PersistentFlags().StringP("query", "", "", "SQL query")
	exportCmd.PersistentFlags().BoolP("schema-only", "", false, "table schema only. no data")
	exportCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")
	exportCmd.PersistentFlags().StringSliceP("hash", "a", []string{"sha256"}, "hash algorithm(sha256 or sha512)")
	exportCmd.PersistentFlags().StringSliceP("enable-target", "", []string{"name", "type"}, "hash target extra item (all or each name)")
	exportCmd.PersistentFlags().StringSliceP("disable-target", "", nil, "hash extra items not to be targeted (all or each name)")
	exportCmd.PersistentFlags().BoolP("sign", "", false, "sign TBLN file")
	exportCmd.PersistentFlags().BoolP("quiet", "q", false, "do not prompt for password.")

	rootCmd.AddCommand(exportCmd)
}

func dbExport(cmd *cobra.Command, args []string) error {
	var err error

	var schema string
	if schema, err = cmd.PersistentFlags().GetString("schema"); err != nil {
		return err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("table"); err != nil {
		return err
	}
	var query string
	if query, err = cmd.PersistentFlags().GetString("query"); err != nil {
		return err
	}
	var schemaOnly bool
	if schemaOnly, err = cmd.PersistentFlags().GetBool("schema-only"); err != nil {
		return err
	}
	if tableName == "" && len(args) >= 1 {
		tableName = args[0]
	}
	if schemaOnly && tableName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires table name, in the case of schema-only option")
	}
	if tableName == "" && query == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires table name or SQL query")
	}
	if srcdbName == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires database driver name")
	}
	var output string
	if output, err = cmd.PersistentFlags().GetString("output"); err != nil {
		return err
	}
	var signF bool
	if signF, err = cmd.PersistentFlags().GetBool("sign"); err != nil {
		return err
	}
	if signF && SecFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires secret key file")
	}
	var quiet bool
	if quiet, err = cmd.PersistentFlags().GetBool("quiet"); err != nil {
		cmd.SilenceUsage = false
		return err
	}

	conn, err := db.Open(srcdbName, srcdsn)
	if err != nil {
		return fmt.Errorf("%s: %s", srcdbName, err)
	}
	defer conn.Close()

	var at *tbln.Tbln
	switch {
	case query != "":
		at, err = db.ReadQueryAll(conn, query)
	case schemaOnly:
		at, err = db.GetTableInfo(conn, schema, tableName)
	default:
		at, err = db.ReadTableAll(conn, schema, tableName)
	}
	if err != nil {
		return err
	}

	err = hash(at, cmd, args)
	if err != nil {
		return err
	}
	if signF {
		privKey, err := key.GetPrivateKey(SecFile, KeyName, !quiet)
		if err != nil {
			return err
		}
		_, err = at.Sign(KeyName, privKey)
		if err != nil {
			return err
		}
	}

	var out *os.File
	if output == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(output)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	return tbln.WriteAll(out, at)
}
