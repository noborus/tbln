package sqlite3

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/noborus/tbln/db"
)

// SQLite3 is dummy struct
type SQLite3 struct{}

func init() {
	driver := SQLite3{}
	db.Register("sqlite3", &driver)
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
func (s *SQLite3) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	return nil, fmt.Errorf("this database is not supported")
}

// GetColumnInfo returns information of a table column as an array.
func (s *SQLite3) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	return nil, fmt.Errorf("this database is not supported")
}
