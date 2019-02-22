package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"

	"github.com/noborus/tbln/db"
)

// MySQL is dummy struct
type MySQL struct {
	mysql.MySQLDriver
}

// Constr is Implement Constraint interface.
type Constr struct {
}

func init() {
	driver := MySQL{}
	constr := Constr{}
	db.Register("mysql", &driver, &constr)
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
func (m *Constr) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	query := `SELECT  column_name
	             FROM information_schema.columns
				WHERE table_schema = database()
				  AND table_name = ?
                  AND column_key = 'PRI'
                ORDER BY ordinal_position;`
	return db.GetPrimaryKey(conn, query, tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (m *Constr) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	query := `SELECT
		        column_default
              , is_nullable
              , data_type AS mysql_type
              , character_maximum_length
              , character_octet_length
              , numeric_precision
              , numeric_scale
              , datetime_precision
	     FROM information_schema.columns
        WHERE table_name = ?
		ORDER BY ordinal_position;`
	return db.GetColumnInfo(conn, query, tableName)
}
