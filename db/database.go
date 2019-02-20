package db

import (
	"database/sql"
	"fmt"
	"regexp"
)

// DBD is sql.DB wrapper.
type DBD struct {
	DB *sql.DB
	Tx *sql.Tx
	Driver
	Name  string
	Ph    string
	Quote string
}

// Open is tbln/db Open.
func Open(name string, dsn string) (*DBD, error) {
	c := Get(name)
	if c == nil {
		c = &NotSupport{}
	}
	db, err := sql.Open(name, dsn)
	dbd := &DBD{
		Name:   name,
		DB:     db,
		Driver: c,
		Ph:     c.PlaceHolder(),
		Quote:  c.Quote(),
	}
	return dbd, err
}

func (db *DBD) quoting(name string) string {
	r := regexp.MustCompile(`[^a-z0-9_]+`)
	if r.MatchString(name) {
		return db.Quote + name + db.Quote
	}
	return name
}

// NotSupport is dummy struct.
type NotSupport struct{}

// GetPrimaryKey is dummy function.
func (n *NotSupport) GetPrimaryKey(db *sql.DB, TableName string) ([]string, error) {
	return nil, fmt.Errorf("this database is not supported")
}

// PlaceHolder returns the placeholer string.
func (n *NotSupport) PlaceHolder() string {
	return "?"
}

// Quote returns the quote string.
func (n *NotSupport) Quote() string {
	return `"`
}
