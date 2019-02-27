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

// GetPrimaryKey returns the primary key as a slice.
func (c *Constr) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	query := `SELECT ccu.column_name
	   	        FROM information_schema.table_constraints tc
	               , information_schema.constraint_column_usage ccu
	           WHERE tc.table_name = $1
	             AND tc.constraint_type = 'PRIMARY KEY'
               AND tc.table_catalog = ccu.table_catalog
						   AND tc.table_schema = ccu.table_schema
						   AND tc.table_name = ccu.table_name
               AND tc.constraint_name = ccu.constraint_name;`
	return db.GetPrimaryKey(conn, query, tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (c *Constr) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	query := `WITH u AS (
		SELECT DISTINCT tc.table_name, cc.column_name, 'YES' as is_unique
			FROM information_schema.table_constraints AS tc
			LEFT JOIN information_schema.constraint_column_usage AS cc
				ON (tc.table_name = cc.table_name AND tc.constraint_type = 'UNIQUE')
		 WHERE tc.table_name = $1
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
	 WHERE t.table_name = $1
	 ORDER BY t.ordinal_position;`
	return db.GetColumnInfo(conn, query, tableName)
}
