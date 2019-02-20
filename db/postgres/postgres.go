package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/noborus/tbln/db"
)

// Postgres is dummy struct
type Postgres struct{}

func init() {
	driver := Postgres{}
	db.Register("postgres", &driver)
}

// GetPrimaryKey returns the primary key as a slice.
func (p *Postgres) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	query :=
		`SELECT ccu.column_name 
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
			return nil, err
		}
		pkeys = append(pkeys, pkey)
	}
	rows.Close()
	return pkeys, nil
}

// PlaceHolder returns the placeholer string.
func (p *Postgres) PlaceHolder() string {
	return "$"
}

// Quote returns the quote string.
func (p *Postgres) Quote() string {
	return `"`
}
