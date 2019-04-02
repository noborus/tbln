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
	var err error
	if schema == "" {
		schema, err = c.GetSchema(conn)
		if err != nil {
			return nil, err
		}
	}
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
	var err error
	if schema == "" {
		schema, err = c.GetSchema(conn)
		if err != nil {
			return nil, err
		}
	}
	query := `SELECT cl.column_default
	 , cl.is_nullable
	 , cl.data_type AS mysql_type
	 , cl.column_type AS mysql_columntype
	 , cl.character_maximum_length
	 , cl.character_octet_length
	 , cl.numeric_precision
	 , cl.numeric_scale
	 , cl.datetime_precision
     , (CASE WHEN tc.constraint_type = 'PRIMARY KEY' THEN 'YES'
             WHEN tc.constraint_type = 'UNIQUE' THEN 'YES'
         END) AS is_unique
  FROM information_schema.columns AS cl
  LEFT JOIN information_schema.key_column_usage AS us
         ON (cl.table_schema = us.table_schema
            AND cl.table_name = us.table_name
            AND cl.column_name = us.column_name)
       LEFT JOIN information_schema.table_constraints AS tc
              ON (us.table_name = tc.table_name
                  AND us.constraint_name = tc.constraint_name
				  AND us.table_schema = tc.table_schema)
	WHERE cl.table_schema = ?
	  AND cl.table_name = ?
	ORDER BY cl.ordinal_position;`
	return db.GetColumnInfo(conn, query, schema, tableName)
}
