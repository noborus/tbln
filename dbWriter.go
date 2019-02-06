package tbln

import (
	"database/sql"
	"fmt"
	"strings"
)

// DBWriter is writer struct.
type DBWriter struct {
	Definition
	db     *sql.DB
	tx     *sql.Tx
	stmt   *sql.Stmt
	Create bool
	Ph     string
}

// NewDBWriter is DB write struct.
func NewDBWriter(db *sql.DB, d Definition, create bool) *DBWriter {
	return &DBWriter{
		Definition: d,
		db:         db,
		Create:     create,
		Ph:         "?",
	}
}

// WriteDefinition is write table definition.
func (tw *DBWriter) WriteDefinition() error {
	if tw.Create {
		err := tw.createTable()
		if err != nil {
			return err
		}
	}
	return tw.prepara()
}

func (tw *DBWriter) createTable() error {
	col := make([]string, len(tw.Names))
	for i := 0; i < len(tw.Names); i++ {
		col[i] = tw.Names[i] + " " + tw.Types[i]
	}
	sql := fmt.Sprintf("CREATE TABLE %s ( %s );",
		tw.name, strings.Join(col, ", "))
	_, err := tw.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (tw *DBWriter) prepara() error {
	var err error
	ph := make([]string, len(tw.Names))
	for i := 0; i < len(tw.Names); i++ {
		if tw.Ph == "$" {
			ph[i] = fmt.Sprintf("$%d", i+1)
		} else {
			ph[i] = fmt.Sprintf("?")
		}
	}
	insert := fmt.Sprintf(
		"INSERT INTO %s ( %s ) VALUES ( %s );",
		tw.name, strings.Join(tw.Names, ", "), strings.Join(ph, ", "))
	tw.stmt, err = tw.db.Prepare(insert)
	return err
}

// WriteRow is write one row.
func (tw *DBWriter) WriteRow(row []string) error {
	r := make([]interface{}, len(row))
	for i, v := range row {
		r[i] = v
	}
	_, err := tw.stmt.Exec(r...)
	return err
}
