// Package postgres contains PostgreSQL specific database dialect of tbln/db.
package postgres

import (
	"database/sql"

	// PostgreSQL driver
	_ "github.com/lib/pq"
	"github.com/noborus/tbln/db"
)

// Constr is Implement Constraint interface.
type Constr struct{}

func init() {
	driver := db.Driver{
		Style:      db.Style{PlaceHolder: "$", Quote: `"`},
		Constraint: &Constr{},
	}
	db.Register("postgres", driver)
}

// GetSchema returns the schema string.
func (c *Constr) GetSchema(conn *sql.DB) (string, error) {
	query := "SELECT current_schema();"
	row := conn.QueryRow(query)
	schema := ""
	if err := row.Scan(&schema); err != nil {
		return "", err
	}
	return schema, nil
}

// GetPrimaryKey returns the primary key as a slice.
func (c *Constr) GetPrimaryKey(conn *sql.DB, schema string, tableName string) ([]string, error) {
	var err error
	if schema == "" {
		schema, err = c.GetSchema(conn)
		if err != nil {
			return nil, err
		}
	}
	query := `SELECT c.column_name
	            FROM information_schema.constraint_column_usage c
				LEFT JOIN information_schema.table_constraints t
				       ON t.constraint_name = c.constraint_name
				  AND t.table_schema = c.table_schema
				  AND t.table_catalog = c.table_catalog
			    WHERE t.table_schema = $1
			      AND t.table_name = $2
				  AND t.constraint_type = 'PRIMARY KEY'`
	return db.GetPrimaryKey(conn, query, schema, tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (c *Constr) GetColumnInfo(conn *sql.DB, schema string, tableName string) (map[string][]interface{}, error) {
	if schema == "" {
		schema = "public"
	}
	query := `WITH u AS (
		SELECT DISTINCT tc.table_name, cc.column_name, 'YES' as is_unique
			FROM information_schema.table_constraints AS tc
			LEFT JOIN information_schema.constraint_column_usage AS cc
				ON (tc.table_name = cc.table_name AND tc.constraint_type = 'UNIQUE')
		 WHERE tc.table_schema = $1
		   AND tc.table_name = $2
	)
	SELECT
			column_default
			, is_nullable
			, data_type AS postgres_type
			, character_maximum_length
  		, character_octet_length
		  , numeric_precision
			, numeric_precision_radix
			, numeric_scale
			, datetime_precision
			, interval_type
			, u.is_unique
		FROM information_schema.columns AS t
		LEFT JOIN u
			ON (t.table_name = u.table_name AND t.column_name = u.column_name)
	 WHERE t.table_schema = $1
	   AND t.table_name = $2
	 ORDER BY t.ordinal_position;`
	return db.GetColumnInfo(conn, query, schema, tableName)
}
