package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/noborus/tbln"
)

// Writer is writer struct.
type Writer struct {
	tbln.Definition
	*DBD
	stmt   *sql.Stmt
	Create bool
}

// NewWriter is DB write struct.
func NewWriter(dbd *DBD, definition tbln.Definition, create bool) *Writer {
	if dbd.Tx == nil {
		var err error
		dbd.Tx, err = dbd.DB.Begin()
		if err != nil {
			return nil
		}
	}
	return &Writer{
		Definition: definition,
		DBD:        dbd,
		Create:     create,
	}
}

// WriteDefinition is create table and insert preparation.
func (tw *Writer) WriteDefinition() error {
	if tw.Names == nil {
		if tw.ColumnNum == 0 {
			return fmt.Errorf("column num is 0")
		}
		tw.Names = make([]string, tw.ColumnNum)
		for i := 0; i < tw.ColumnNum; i++ {
			tw.Names[i] = fmt.Sprintf("c%d", i+1)
		}
	}
	if tw.Types == nil {
		if tw.ColumnNum == 0 {
			return fmt.Errorf("column num is 0")
		}
		tw.Types = make([]string, tw.ColumnNum)
		for i := 0; i < tw.ColumnNum; i++ {
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

func (tw *Writer) createTable() error {
	col := make([]string, len(tw.Names))
	for i := 0; i < len(tw.Names); i++ {
		col[i] = tw.quoting(tw.Names[i]) + " " + tw.Types[i]
	}
	sql := fmt.Sprintf("CREATE TABLE %s ( %s );",
		tw.quoting(tw.TableName), strings.Join(col, ", "))
	_, err := tw.Tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("%s: %s", err, sql)
	}
	return nil
}

func (tw *Writer) prepara() error {
	var err error
	names := make([]string, len(tw.Names))
	ph := make([]string, len(tw.Names))
	for i := 0; i < len(tw.Names); i++ {
		names[i] = tw.quoting(tw.Names[i])
		if tw.Ph == "$" {
			ph[i] = fmt.Sprintf("$%d", i+1)
		} else {
			ph[i] = "?"
		}
	}
	insert := fmt.Sprintf(
		"INSERT INTO %s ( %s ) VALUES ( %s );",
		tw.quoting(tw.TableName), strings.Join(names, ", "), strings.Join(ph, ", "))
	tw.stmt, err = tw.Tx.Prepare(insert)
	if err != nil {
		return fmt.Errorf("%s: %s", err, insert)
	}
	return nil
}

func (tw *Writer) convertDBType(dbtype string, value string) interface{} {
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
func (tw *Writer) WriteRow(row []string) error {
	r := make([]interface{}, len(row))
	for i, v := range row {
		r[i] = tw.convertDBType(tw.Types[i], v)
	}
	_, err := tw.stmt.Exec(r...)
	return err
}

// WriteTable writes all rows to the table.
func WriteTable(dbd *DBD, tbln *tbln.Tbln, create bool) error {
	w := NewWriter(dbd, tbln.Definition, create)
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
