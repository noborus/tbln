// Package mysql contains MySQL specific database dialect of tbln/db.
package mysql

import (
	"database/sql"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/noborus/tbln/db"
)

// Constr is Implement Constraint interface.
type Constr struct{}

func init() {
	driver := db.Driver{
		Style:      db.Style{PlaceHolder: "?", Quote: "`"},
		Constraint: &Constr{},
	}
	db.Register("mysql", driver)
}

// GetSchema returns the schema string.
func (c *Constr) GetSchema(conn *sql.DB) (string, error) {
	query := "SELECT schema();"
	row := conn.QueryRow(query)
	schema := ""
	if err := row.Scan(&schema); err != nil {
		return "", err
	}
	return schema, nil
}

// GetPrimaryKey returns the primary key as a slice.
func (c *Constr) GetPrimaryKey(conn *sql.DB, schema string, tableName string) ([]string, error) {
	query := `SELECT  column_name
	            FROM information_schema.columns
				     WHERE table_schema = ?
				       AND table_name = ?
               AND column_key = 'PRI'
			 ORDER BY ordinal_position;`
	return db.GetPrimaryKey(conn, query, schema, tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (c *Constr) GetColumnInfo(conn *sql.DB, schema string, tableName string) (map[string][]interface{}, error) {
	query := `WITH u AS (
		SELECT DISTINCT tc.table_name, cc.column_name, 'YES' as constraint_unique
			FROM information_schema.table_constraints AS tc
			LEFT JOIN information_schema.KEY_COLUMN_USAGE AS cc
				ON (tc.table_name = cc.table_name AND tc.constraint_type = 'UNIQUE')
		 WHERE tc.table_schema = ?
		 AND tc.table_name = ?
	)
	   SELECT column_default
			  , is_nullable
			  , data_type AS mysql_type
			  , character_maximum_length
			  , character_octet_length
			  , numeric_precision
			  , numeric_scale
			  , datetime_precision
			  , u.constraint_unique
			FROM information_schema.columns AS tc
			LEFT JOIN u
				     ON (tc.table_name = u.table_name
				    AND tc.column_name = u.column_name )
		WHERE tc.table_schema = ?
	      AND tc.table_name = ?
		 ORDER BY ordinal_position;`
	return db.GetColumnInfo(conn, query, schema, tableName, schema, tableName)
}
