package tbln

import (
	"database/sql"
	"strings"
)

// DBWriter is writer struct.
type DBWriter struct {
	Table
	db      *sql.DB
	tx      *sql.Tx
	preSQL  string
	postSQL string
}

// NewDBWriter is DB write struct.
func NewDBWriter(db *sql.DB, tbl Table) *DBWriter {
	return &DBWriter{
		Table: tbl,
		db:    db,
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
	tw.preSQL = "INSERT INTO " + tw.name + " ("
	tw.preSQL += strings.Join(tw.Names, ", ")
	tw.preSQL += ") VALUES ("
	tw.postSQL = ");"
	return nil
}

// WriteRow is write one row.
func (tw *DBWriter) WriteRow(row []string) error {
	sql := tw.preSQL + "'" + strings.Join(row, "', '") + "'" + tw.postSQL
	_, err := tw.db.Exec(sql)
	return err
}
