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
	sdb, err := sql.Open(driver, dsn)
	ph := "?"
	quote := "`"
	if driver == "postgres" {
		ph = "$"
		quote = `"`
	}
	dbd := &DBD{
		Name:  driver,
		DB:    sdb,
		Ph:    ph,
		Quote: quote,
	}
	return dbd, err
}
