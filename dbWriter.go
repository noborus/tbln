package tbln

import (
	"database/sql"
	"fmt"
	"strings"
)

// DBWriter is writer struct.
type DBWriter struct {
	Table
	db   *sql.DB
	tx   *sql.Tx
	stmt *sql.Stmt
	Ph   string
}

// NewDBWriter is DB write struct.
func NewDBWriter(db *sql.DB, tbl Table) *DBWriter {
	return &DBWriter{
		Table: tbl,
		db:    db,
		Ph:    "?",
	}
}

// WriteInfo is output table information.
func (tw *DBWriter) WriteInfo() error {
	sql := "CREATE TABLE " + tw.name + " ("
	col := make([]string, len(tw.Names))
	for i := 0; i < len(tw.Names); i++ {
		col[i] = tw.Names[i] + " " + tw.Types[i]
	}
	sql += strings.Join(col, ", ")
	sql += ");"
	_, err := tw.db.Exec(sql)
	if err != nil {
		return err
	}
	tw.preparation()
	return nil
}

func (tw *DBWriter) preparation() error {
	var err error
	ph := make([]string, len(tw.Names))
	for i := 0; i < len(tw.Names); i++ {
		if tw.Ph == "$" {
			ph[i] = fmt.Sprintf("$%d", i+1)
		} else {
			ph[i] = fmt.Sprintf("?")
		}
	}
	insert := "INSERT INTO " + tw.name + " ("
	insert += strings.Join(tw.Names, ", ")
	insert += ") VALUES ("
	insert += strings.Join(ph, ", ")
	insert += ");"
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
