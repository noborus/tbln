package tbln

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/noborus/tbln/db"
)

type NotSupport struct{}

func (n *NotSupport) GetPrimaryKey(db *sql.DB, tableName string) ([]string, error) {
	return nil, fmt.Errorf("this database is not supported")
}

// DBD is sql.DB wrapper.
type DBD struct {
	DB *sql.DB
	Tx *sql.Tx
	db.Constraint
	Name  string
	Ph    string
	Quote string
}

// DBOpen is *sql.Open Wrapper.
func DBOpen(driver string, dsn string) (*DBD, error) {
	var c db.Constraint
	if db.DB[driver] != nil {
		c = db.DB[driver]
	} else {
		c = &NotSupport{}
	}

	db, err := sql.Open(driver, dsn)
	var ph, quote string
	switch driver {
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
		Name:       driver,
		DB:         db,
		Constraint: c,
		Ph:         ph,
		Quote:      quote,
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
