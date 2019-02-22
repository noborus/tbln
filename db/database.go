package db

import (
	"database/sql"
	"regexp"
)

// TDB is sql.DB wrapper.
type TDB struct {
	DB *sql.DB
	Tx *sql.Tx
	Style
	Constraint
	Name string
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

func (TDB *TDB) quoting(name string) string {
	r := regexp.MustCompile(`[^a-z\.0-9_]+`)
	q := TDB.Style.Quote
	escape := regexp.MustCompile(`(` + q + `)`)
	if r.MatchString(name) {
		name = escape.ReplaceAllString(name, "$1"+TDB.Style.Quote)
		return q + name + q
	}
	return name
}

// GetPrimaryKey returns the primary key as a slice.
func GetPrimaryKey(conn *sql.DB, query string, tableName string) ([]string, error) {
	rows, err := conn.Query(query, tableName)
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
	rows.Close()
	return pkeys, nil
}

// GetColumnInfo returns information of a table column as an array.
func GetColumnInfo(conn *sql.DB, query string, tableName string) (map[string][]interface{}, error) {
	rows, err := conn.Query(query, tableName)
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
