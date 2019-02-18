package tbln

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// DBWriter is writer struct.
type DBWriter struct {
	Definition
	*DBD
	db     *sql.DB
	tx     *sql.Tx
	stmt   *sql.Stmt
	Create bool
}

// NewDBWriter is DB write struct.
func NewDBWriter(dbd *DBD, definition Definition, create bool) *DBWriter {
	if dbd.Tx == nil {
		var err error
		dbd.Tx, err = dbd.DB.Begin()
		if err != nil {
			return nil
		}
	}
	return &DBWriter{
		Definition: definition,
		DBD:        dbd,
		db:         dbd.DB,
		tx:         dbd.Tx,
		Create:     create,
	}
}

// WriteDefinition is create table and insert preparation.
func (tw *DBWriter) WriteDefinition() error {
	if tw.Names == nil {
		if tw.columnNum == 0 {
			return fmt.Errorf("column num is 0")
		}
		tw.Names = make([]string, tw.columnNum)
		for i := 0; i < tw.columnNum; i++ {
			tw.Names[i] = fmt.Sprintf("c%d", i+1)
		}
	}
	if tw.Types == nil {
		if tw.columnNum == 0 {
			return fmt.Errorf("column num is 0")
		}
		tw.Types = make([]string, tw.columnNum)
		for i := 0; i < tw.columnNum; i++ {
			tw.Types[i] = "text"
		}
	}
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
		col[i] = tw.quoting(tw.Names[i]) + " " + tw.Types[i]
	}
	sql := fmt.Sprintf("CREATE TABLE %s ( %s );",
		tw.quoting(tw.tableName), strings.Join(col, ", "))
	_, err := tw.tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("%s: %s", err, sql)
	}
	return nil
}

func (tw *DBWriter) prepara() error {
	var err error
	names := make([]string, len(tw.Names))
	ph := make([]string, len(tw.Names))
	for i := 0; i < len(tw.Names); i++ {
		names[i] = tw.quoting(tw.Names[i])
		if tw.Ph == "$" {
			ph[i] = fmt.Sprintf("$%d", i+1)
		} else {
			ph[i] = fmt.Sprintf("?")
		}
	}
	insert := fmt.Sprintf(
		"INSERT INTO %s ( %s ) VALUES ( %s );",
		tw.quoting(tw.tableName), strings.Join(names, ", "), strings.Join(ph, ", "))
	tw.stmt, err = tw.tx.Prepare(insert)
	if err != nil {
		return fmt.Errorf("%s: %s", err, insert)
	}
	return nil
}

func (tw *DBWriter) convertDBType(dbtype string, value string) interface{} {
	switch strings.ToLower(dbtype) {
	case "datetime", "timestamp":
		if tw.DBD.Name == "mysql" {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil
			}
			return t
		}
		return value
	}
	return value
}

// WriteRow is write one row.
func (tw *DBWriter) WriteRow(row []string) error {
	r := make([]interface{}, len(row))
	for i, v := range row {
		r[i] = tw.convertDBType(tw.Types[i], v)
	}
	_, err := tw.stmt.Exec(r...)
	return err
}

// WriteTable writes all rows to the table.
func WriteTable(db *DBD, tbln *Tbln, create bool) error {
	w := NewDBWriter(db, tbln.Definition, create)
	err := w.WriteDefinition()
	if err != nil {
		return err
	}
	for _, row := range tbln.Rows {
		err = w.WriteRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}
