package tbln

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/noborus/tbln/driver/postgres"
)

type NotSupport struct{}

// Constraint is get the Constraint for each database
type Constraint interface {
	GetPrimaryKey(db *sql.DB, tableName string) ([]string, error)
}

func (n *NotSupport) GetPrimaryKey(db *sql.DB, tableName string) ([]string, error) {
	return nil, fmt.Errorf("this database is not supported")
}

// DBD is sql.DB wrapper.
type DBD struct {
	DB *sql.DB
	Tx *sql.Tx
	Constraint
	Name  string
	Ph    string
	Quote string
}

// DBOpen is *sql.Open Wrapper.
func DBOpen(driver string, dsn string) (*DBD, error) {
	db, err := sql.Open(driver, dsn)
	var c Constraint
	c = &NotSupport{}
	var ph, quote string
	switch driver {
	case "postgres":
		ph = "$"
		quote = `"`
		c = postgres.Postgres{}
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
