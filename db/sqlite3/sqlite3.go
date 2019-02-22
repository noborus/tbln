package sqlite3

import (
	"database/sql"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/noborus/tbln/db"
)

// SQLite3 struct
type SQLite3 struct {
	sqlite3.SQLiteDriver
}

// Constr is Implement Constraint interface.
type Constr struct {
}

func init() {
	driver := SQLite3{}
	constr := Constr{}

	db.Register("sqlite3", &driver, &constr)
}

// PlaceHolder returns the placeholer string.
func (s *SQLite3) PlaceHolder() string {
	return "?"
}

// Quote returns the quote string.
func (s *SQLite3) Quote() string {
	return "`"
}

// GetPrimaryKey returns the primary key as a slice.
func (p *Constr) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	query := `SELECT name FROM PRAGMA_TABLE_INFO(?)`
	return db.GetPrimaryKey(conn, query, tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (p *Constr) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	query := `SELECT
	         type AS sqlite3_type
				 , CASE  "notnull"
				     WHEN 0 THEN 'YES'
					 ELSE 'NO'
				    END AS is_nullable
				FROM PRAGMA_TABLE_INFO(?)`
	return db.GetColumnInfo(conn, query, tableName)
}
