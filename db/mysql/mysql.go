package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/noborus/tbln/db"
)

// MySQL is dummy struct
type MySQL struct{}

func init() {
	driver := MySQL{}
	db.Register("mysql", &driver)
}

// PlaceHolder returns the placeholer string.
func (m *MySQL) PlaceHolder() string {
	return "?"
}

// Quote returns the quote string.
func (m *MySQL) Quote() string {
	return "`"
}

// GetPrimaryKey returns the primary key as a slice.
func (m *MySQL) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	return nil, fmt.Errorf("this database is not supported")
}
