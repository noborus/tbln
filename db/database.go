// Package db translates tbln and database tables.
//
// Read and write the table using *sql.DB.
// PostgreSQL, MySQL and SQLlite3 interprets SQL dialects
// and outputs more detailed information.
package db

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
)

// TDB struct is sql.DB wrapper.
type TDB struct {
	Name string
	*sql.DB
	*sql.Tx
	IsTx bool
	Style
	Constraint
}

// Open is tbln/db Open.
func Open(name string, dsn string) (*TDB, error) {
	d := GetDriver(name)
	if d.PlaceHolder == "" {
		// Default SQL Style.
		d.Style = Style{
			PlaceHolder: "?",
			Quote:       `"`,
		}
	}
	if d.Constraint == nil {
		d.Constraint = &DefaultConstr{}
	}
	db, err := sql.Open(name, dsn)
	TDB := &TDB{
		Name:       name,
		DB:         db,
		Style:      d.Style,
		Constraint: d.Constraint,
	}
	return TDB, err
}

// Begin is the sql.DB Begin wrapper.
func (tdb *TDB) Begin() error {
	var err error
	if tdb.IsTx {
		return fmt.Errorf("already in transaction")
	}
	tdb.Tx, err = tdb.DB.Begin()
	if err != nil {
		return err
	}
	tdb.IsTx = true
	return nil
}

// Commit is the sql.DB Commit wrapper.
func (tdb *TDB) Commit() error {
	if !tdb.IsTx {
		if tdb.Tx != nil {
			tdb.Tx.Rollback()
		}
		return fmt.Errorf("no transaction")
	}
	tdb.IsTx = false
	err := tdb.Tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (tdb *TDB) quoting(obj string) string {
	r := regexp.MustCompile(`[^a-z0-9_]+`)
	q := tdb.Style.Quote
	escape := regexp.MustCompile(`(` + q + `)`)
	if r.MatchString(obj) {
		obj = escape.ReplaceAllString(obj, "$1"+tdb.Style.Quote)
		return q + obj + q
	}
	return obj
}

// GetPrimaryKey returns the primary key as a slice.
func GetPrimaryKey(conn *sql.DB, query string, schema string, tableName string) ([]string, error) {
	debug.Printf("SQL:%s", query)
	rows, err := conn.Query(query, schema, tableName)
	if err != nil {
		return nil, err
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
	return pkeys, rows.Close()

}

// GetColumnInfo returns information of a table column as an array.
func GetColumnInfo(conn *sql.DB, query string, args ...interface{}) (map[string][]interface{}, error) {
	rows, err := conn.Query(query, args...)
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

var debug = debugT(false)

type debugT bool

func (d debugT) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}

// Debug changes Debug output.
func Debug(d debugT) {
	debug = d
}
