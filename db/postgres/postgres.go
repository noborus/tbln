package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/noborus/tbln/db"
)

// Postgres is dummy struct
type Postgres struct {
	pq.Driver
}

func init() {
	driver := Postgres{}
	constr := Constr{}

	db.Register("postgres", &driver, &constr)
}

// PlaceHolder returns the placeholer string.
func (p *Postgres) PlaceHolder() string {
	return "$"
}

// Quote returns the quote string.
func (p *Postgres) Quote() string {
	return `"`
}

// Constr is Implement Constraint interface.
type Constr struct {
}

// GetPrimaryKey returns the primary key as a slice.
func (p *Constr) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	query := `SELECT ccu.column_name
	   	      FROM information_schema.table_constraints tc
	             , information_schema.constraint_column_usage ccu
	         WHERE tc.table_name = $1
	           AND tc.constraint_type = 'PRIMARY KEY'
             AND tc.table_catalog = ccu.table_catalog
						 AND tc.table_schema = ccu.table_schema
						 AND tc.table_name = ccu.table_name
					   AND tc.constraint_name = ccu.constraint_name;`
	return db.GetPrimaryKey(conn, query, tableName)
}

// GetColumnInfo returns information of a table column as an array.
func (p *Constr) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	query := `SELECT column_name
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
	return db.GetColumnInfo(conn, query, tableName)
}
