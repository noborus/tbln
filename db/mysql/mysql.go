package mysql

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/noborus/tbln/db"
)

// MySQL is dummy struct
type MySQL struct {
	mysql.MySQLDriver
}

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
	return nil, fmt.Errorf("thiis database is not supported")
}

// columnInfoQuery is a query to get ColumnInfo.
var columnInfoQuery = `SELECT column_name
		      , column_default
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

// GetColumnInfo returns information of a table column as an array.
func (m *MySQL) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	rows, err := conn.Query(columnInfoQuery, tableName)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	info := make(map[string][]interface{})
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		for i, name := range columns {
			col := info[name]
			col = append(col, values[i])
			info[name] = col
		}
	}
	return info, nil
}
