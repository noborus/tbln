package sqlite3

import (
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/noborus/tbln/db"
)

// SQLite3 is dummy struct
type SQLite3 struct {
	sqlite3.SQLiteDriver
}

func init() {
	driver := SQLite3{}
	db.Register("sqlite3", &driver, nil)
}

// PlaceHolder returns the placeholer string.
func (s *SQLite3) PlaceHolder() string {
	return "?"
}

// Quote returns the quote string.
func (s *SQLite3) Quote() string {
	return "`"
}
