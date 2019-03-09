// Package sqlite3 contains SQLite3 specific database dialect of tbln/db.
package sqlite3

import (
	"database/sql"

	// SQLite3 driver.
	_ "github.com/mattn/go-sqlite3"
	"github.com/noborus/tbln/db"
)

// Constr is Implement Constraint interface.
type Constr struct{}

func init() {
	driver := db.Driver{
		Style:      db.Style{PlaceHolder: "?", Quote: "`"},
		Constraint: &Constr{},
	}
	db.Register("sqlite3", driver)
}

// GetSchema returns the schema string.
func (c *Constr) GetSchema(conn *sql.DB) (string, error) {
	return "", nil
}

// GetPrimaryKey returns the primary key as a slice.
func (c *Constr) GetPrimaryKey(conn *sql.DB, schema string, tableName string) ([]string, error) {
	query := `SELECT name FROM PRAGMA_TABLE_INFO(?) WHERE pk = 1`
	return db.GetPrimaryKey(conn, query, "", tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (c *Constr) GetColumnInfo(conn *sql.DB, schema string, tableName string) (map[string][]interface{}, error) {
	query := `SELECT
					  t.type AS sqlite3_type
					  , t.dflt_value AS column_default
					  , CASE  "notnull"
					    WHEN 0 THEN 'YES'
                        ELSE 'NO'
						END AS is_nullable
					  , CASE 
						WHEN i.cid IS NOT NULL THEN 'YES'
						END AS is_unique
				FROM pragma_table_info(?) AS t 
				LEFT JOIN (pragma_index_list(?) AS li
					 LEFT JOIN pragma_index_info(li.name) AS il) AS i
					   ON t.cid = i.cid
			   ORDER BY t.cid;`
	return db.GetColumnInfo(conn, query, tableName, tableName)
}
