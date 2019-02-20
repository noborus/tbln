package db

import (
	"database/sql"
	"regexp"
)

// TDB is sql.DB wrapper.
type TDB struct {
	DB *sql.DB
	Tx *sql.Tx
	Driver
	Name string
}

// Open is tbln/db Open.
func Open(name string, dsn string) (*TDB, error) {
	d := GetDriver(name)
	if d == nil {
		d = &Default{}
	}
	db, err := sql.Open(name, dsn)
	TDB := &TDB{
		Name:   name,
		DB:     db,
		Driver: d,
	}
	return TDB, err
}

func (TDB *TDB) quoting(name string) string {
	r := regexp.MustCompile(`[^a-z0-9_]+`)
	q := TDB.Driver.Quote()
	escape := regexp.MustCompile(`(` + q + `)`)
	if r.MatchString(name) {
		name = escape.ReplaceAllString(name, "$1"+TDB.Driver.Quote())
		return q + name + q
	}
	return name
}
