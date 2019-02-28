package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/noborus/tbln"
)

type CreateMode int

const (
	NotCreate = iota
	Create
	IfNotExists
	ReCreate
)

// Writer writes records to database table.
type Writer struct {
	*tbln.Definition
	*TDB
	stmt  *sql.Stmt
	cmode CreateMode
}

// NewWriter returns a new Writer that writes to database table.
func NewWriter(tdb *TDB, definition *tbln.Definition, cmode CreateMode) *Writer {
	if tdb.Tx == nil {
		var err error
		tdb.Tx, err = tdb.Begin()
		if err != nil {
			return nil
		}
	}
	return &Writer{
		Definition: definition,
		TDB:        tdb,
		cmode:      cmode,
	}
}

// WriteRow writes a single tbln record to w.
// A record is a slice of strings with each string being one field.
func (w *Writer) WriteRow(row []string) error {
	r := make([]interface{}, len(row))
	for i, v := range row {
		r[i] = w.convertDBType(w.Types[i], v)
	}
	_, err := w.stmt.Exec(r...)
	return err
}

// WriteTable writes all rows to the table.
func WriteTable(tdb *TDB, tbln *tbln.Tbln, cmode CreateMode) error {
	w := NewWriter(tdb, tbln.Definition, cmode)
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

// WriteDefinition is create table and insert preparation.
func (w *Writer) WriteDefinition() error {
	if w.Names == nil {
		if w.ColumnNum() == 0 {
			return fmt.Errorf("column num is 0")
		}
		w.Names = make([]string, w.ColumnNum())
		for i := 0; i < w.ColumnNum(); i++ {
			w.Names[i] = fmt.Sprintf("c%d", i+1)
		}
	}
	if w.Types == nil {
		if w.ColumnNum() == 0 {
			return fmt.Errorf("column num is 0")
		}
		w.Types = make([]string, w.ColumnNum())
		for i := 0; i < w.ColumnNum(); i++ {
			w.Types[i] = "text"
		}
	}
	if w.cmode > NotCreate {
		err := w.createTable()
		if err != nil {
			return err
		}
	}
	return w.prepara()
}

func (w *Writer) dropTable() error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", w.quoting(w.TableName()))
	_, err := w.Tx.Exec(sql)
	fmt.Println(sql)
	if err != nil {
		return fmt.Errorf("%s: %s", err, sql)
	}
	return err
}

func (w *Writer) createTable() error {
	mode := ""
	if w.cmode == ReCreate {
		err := w.dropTable()
		if err != nil {
			return err
		}
	} else if w.cmode == IfNotExists {
		mode = "IF NOT EXISTS"
	}
	col := make([]string, len(w.Names))
	for i := 0; i < len(w.Names); i++ {
		col[i] = w.quoting(w.Names[i]) + " " + w.Types[i]
	}
	sql := fmt.Sprintf("CREATE TABLE %s %s ( %s );",
		mode, w.quoting(w.TableName()), strings.Join(col, ", "))
	fmt.Println(sql)
	_, err := w.Tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("%s: %s", err, sql)
	}
	return nil
}

func (w *Writer) prepara() error {
	var err error
	names := make([]string, len(w.Names))
	ph := make([]string, len(w.Names))
	for i := 0; i < len(w.Names); i++ {
		names[i] = w.quoting(w.Names[i])
		if w.Style.PlaceHolder == "$" {
			ph[i] = fmt.Sprintf("$%d", i+1)
		} else {
			ph[i] = "?"
		}
	}
	insert := fmt.Sprintf(
		"INSERT INTO %s ( %s ) VALUES ( %s );",
		w.quoting(w.TableName()), strings.Join(names, ", "), strings.Join(ph, ", "))
	w.stmt, err = w.Tx.Prepare(insert)
	if err != nil {
		return fmt.Errorf("%s: %s", err, insert)
	}
	return nil
}

func (w *Writer) convertDBType(dbtype string, value string) interface{} {
	switch strings.ToLower(dbtype) {
	case "datetime", "timestamp":
		if w.TDB.Name == "mysql" {
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
