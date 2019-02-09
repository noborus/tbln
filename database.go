package tbln

import "database/sql"

// DBD is sql.DB wrapper.
type DBD struct {
	DB    *sql.DB
	Name  string
	Ph    string
	Quote string
}

// DBOpen is *sql.Open Wrapper.
func DBOpen(driver string, dsn string) (*DBD, error) {
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
		Name:  driver,
		DB:    db,
		Ph:    ph,
		Quote: quote,
	}
	return dbd, err
}
