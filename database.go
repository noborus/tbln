package tbln

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/noborus/tbln/db"
)

// DBD is sql.DB wrapper.
type DBD struct {
	DB *sql.DB
	Tx *sql.Tx
	db.Driver
	Name  string
	Ph    string
	Quote string
}

// DBOpen is *sql.Open Wrapper.
func DBOpen(name string, dsn string) (*DBD, error) {
	var c db.Driver
	if c = db.Get(name); c == nil {
		c = &NotSupport{}
	}

	db, err := sql.Open(name, dsn)
	var ph, quote string
	switch name {
	case "postgres":
		ph = "$"
		quote = `"`
	case "mysql":
		ph = "?"
		quote = "`"
	case "sqlite3":
		ph = "?"
		quote = "`"
	default:
		// SQL standard
		ph = "?"
		quote = `"`
	}
	dbd := &DBD{
		Name:   name,
		DB:     db,
		Driver: c,
		Ph:     ph,
		Quote:  quote,
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
func (n *NotSupport) GetPrimaryKey(db *sql.DB, tableName string) ([]string, error) {
	return nil, fmt.Errorf("this database is not supported")
}
