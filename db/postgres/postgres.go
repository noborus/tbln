package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/noborus/tbln/db"
)

// Postgres is dummy struct
type Postgres struct {
	pq.Driver
}

func init() {
	driver := Postgres{}
	db.Register("postgres", &driver)
}

// PlaceHolder returns the placeholer string.
func (p *Postgres) PlaceHolder() string {
	return "$"
}

// Quote returns the quote string.
func (p *Postgres) Quote() string {
	return `"`
}

// GetPrimaryKey returns the primary key as a slice.
func (p *Postgres) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	query :=
		`        SELECT ccu.column_name
	   	         FROM information_schema.table_constraints tc
	                 , information_schema.constraint_column_usage ccu
	             WHERE tc.table_name = $1
	               AND tc.constraint_type = 'PRIMARY KEY'
	               AND tc.table_catalog = ccu.table_catalog
	               AND tc.table_schema = ccu.table_schema
	               AND tc.table_name = ccu.table_name
	   					AND tc.constraint_name = ccu.constraint_name;`
	rows, err := conn.Query(query, tableName)
	if err != nil {
		return nil, err.(*pq.Error)
	}
	var pkeys []string
	for rows.Next() {
		var pkey string
		err = rows.Scan(&pkey)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		pkeys = append(pkeys, pkey)
	}
	rows.Close()
	return pkeys, nil
}

// columnInfoQuery is a query to get ColumnInfo.
var columnInfoQuery = `SELECT column_name
		      , column_default
              , is_nullable
              , data_type AS postgres_type
              , character_maximum_length
              , character_octet_length
              , numeric_precision
              , numeric_precision_radix
              , numeric_scale
              , datetime_precision
              , interval_type
	     FROM information_schema.columns
        WHERE table_catalog = current_database()
		    AND table_name = $1
	    ORDER BY ordinal_position;`

// GetColumnInfo returns information of a table column as an array.
func (p *Postgres) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	rows, err := conn.Query(columnInfoQuery, tableName)
	if err != nil {
		return nil, err.(*pq.Error)
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
