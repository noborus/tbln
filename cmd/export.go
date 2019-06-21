package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xo/dburl"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"

	// MySQL driver
	_ "github.com/noborus/tbln/db/mysql"
	// PostgreSQL driver
	_ "github.com/noborus/tbln/db/postgres"
	// SQLite3 driver
	_ "github.com/noborus/tbln/db/sqlite3"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:          "export [flags] <Table Name> or <SQL Query>",
	SilenceUsage: true,
	Short:        "Export database table or query",
	Long: `Export from the database by table or SQL Query.
Add as much of the table information as possible including 
the column name and column type.`,
	RunE: cmdExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().StringP("dburl", "d", "", "database url")
	exportCmd.PersistentFlags().StringP("schema", "n", "", "schema name")
	exportCmd.PersistentFlags().StringP("table", "t", "", "table name")
	exportCmd.PersistentFlags().StringP("query", "x", "", "SQL query")
	exportCmd.PersistentFlags().BoolP("schema-only", "", false, "table schema only. no data")
	exportCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")
	exportCmd.PersistentFlags().StringSliceP("hash", "a", []string{"sha256"}, "hash algorithm(sha256 or sha512)")
	exportCmd.PersistentFlags().StringSliceP("enable-target", "", []string{"name", "type"}, "hash target extra item (all or each name)")
	exportCmd.PersistentFlags().StringSliceP("disable-target", "", nil, "hash extra items not to be targeted (all or each name)")
	exportCmd.PersistentFlags().BoolP("sign", "", false, "sign TBLN file")
	exportCmd.PersistentFlags().BoolP("quiet", "q", false, "do not prompt for password.")
}

func cmdExport(cmd *cobra.Command, args []string) error {
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

	var signF bool
	if signF, err = cmd.PersistentFlags().GetBool("sign"); err != nil {
		return err
	}
	if signF && SecFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires secret key file")
	}

	var url string
	if url, err = cmd.PersistentFlags().GetString("dburl"); err != nil {
		return err
	}
	if url == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("database url not found")
	}
	u, err := dburl.Parse(url)
	if err != nil {
		return fmt.Errorf("%s: %s", url, err)
	}

	tb, err := dbExport(u, query, schema, tableName, schemaOnly)
	if err != nil {
		return err
	}

	tb, err = hashFile(tb, cmd)
	if err != nil {
		return err
	}
	if signF {
		tb, err = signFile(tb, cmd)
		if err != nil {
			return err
		}
	}

	return outputFile(tb, cmd)
}

func dbExport(u *dburl.URL, query string, schema string, tableName string, schemaOnly bool) (*tbln.TBLN, error) {
	conn, err := db.Open(u.Driver, u.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", u.Driver, err)
	}
	defer conn.Close()
	switch {
	case query != "":
		return db.ReadQueryAll(conn, query)
	case schemaOnly:
		return db.GetTableInfo(conn, schema, tableName)
	default:
		return db.ReadTableAll(conn, schema, tableName)
	}
}
