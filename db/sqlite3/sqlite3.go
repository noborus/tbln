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

// GetPrimaryKey returns the primary key as a slice.
func (p *Constr) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	query := `SELECT name FROM PRAGMA_TABLE_INFO(?) WHERE pk = 1`
	return db.GetPrimaryKey(conn, query, tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (p *Constr) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	query := `SELECT type AS sqlite3_type
                 , dflt_value AS column_default
                 , CASE  "notnull"
                     WHEN 0 THEN 'YES'
                     ELSE 'NO'
                     END AS is_nullable
              FROM PRAGMA_TABLE_INFO(?)`
	return db.GetColumnInfo(conn, query, tableName)
}
