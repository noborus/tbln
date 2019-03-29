package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/noborus/tbln"
)

// CreateMode represents the mode of table creation.
type CreateMode int

const (
	// NotCreate does not execute CREATE TABLE
	NotCreate = iota
	// Create is mormal creation.
	Create
	// IfNotExists does nothing if already exists.
	IfNotExists
	// ReCreate erases the table and creates it again.
	ReCreate
	// CreateOnly does only CREATE TABLE.
	CreateOnly
)

// InsertMode represents the mode of insert conflicts.
type InsertMode int

const (
	// Normal does normal insert.
	Normal = iota
	// OrIgnore  ignores at insert conflict.
	OrIgnore
)

// Writer writes records to database table.
type Writer struct {
	*tbln.Definition
	*TDB
	tableFullName string // schema.table
	stmt          *sql.Stmt
	cmode         CreateMode
	imode         InsertMode
	ReplaceLN     bool
}

// NewWriter returns a new Writer that writes to database table.
func NewWriter(tdb *TDB, definition *tbln.Definition) (*Writer, error) {
	if tdb.Tx == nil {
		err := tdb.Begin()
		if err != nil {
			return nil, err
		}
	}
	return &Writer{
		Definition: definition,
		TDB:        tdb,
		ReplaceLN:  true,
	}, nil
}

// WriteRow writes a single tbln record to w.
// A record is a slice of strings with each string being one field.
func (w *Writer) WriteRow(row []string) error {
	r := make([]interface{}, len(row))
	for i, v := range row {
		if w.ReplaceLN {
			v = strings.ReplaceAll(v, "\\n", "\n")
		}
		r[i] = w.convertDBType(w.Types[i], v)
	}
	_, err := w.stmt.Exec(r...)
	return err
}

// WriteTable writes all rows to the table.
func WriteTable(tdb *TDB, tbln *tbln.Tbln, schema string, cmode CreateMode, imode InsertMode) error {
	w, err := NewWriter(tdb, tbln.Definition)
	if err != nil {
		return err
	}
	w.cmode = cmode
	w.imode = imode
	if schema != "" {
		w.tableFullName = w.quoting(schema) + "." + w.quoting(w.TableName())
	} else {
		w.tableFullName = w.quoting(w.TableName())
	}
	err = w.WriteDefinition()
	if err != nil {
		return err
	}
	if cmode == CreateOnly {
		return nil
	}
	for _, row := range tbln.Rows {
		err = w.WriteRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteDefinition is create table and insert prepare.
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
	return w.prepare()
}

func (w *Writer) dropTable() error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", w.tableFullName)
	debug.Printf("SQL:%s", sql)
	_, err := w.Tx.Exec(sql)
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
		mode = "IF NOT EXISTS "
	}
	constraints := w.createConstraints()
	driverName := w.TDB.Name
	var typeNames []string
	if _, ok := w.Extras[driverName+"_type"]; ok {
		var charMax []string
		if _, ok := w.Extras["character_maximum_length"]; ok {
			charMax = tbln.SplitRow(toString(w.ExtraValue("character_maximum_length")))
		}
		typeNames = tbln.SplitRow(toString(w.ExtraValue(driverName + "_type")))
		for i, cm := range charMax {
			if cm != "" {
				typeNames[i] = typeNames[i] + "(" + cm + ")"
			}
		}
	}
	if len(typeNames) == 0 {
		typeNames = w.Types
	}
	col := make([]string, len(w.Names))
	for i := 0; i < len(w.Names); i++ {
		col[i] = w.quoting(w.Names[i]) + " " + typeNames[i] + constraints[i]
	}
	sql := fmt.Sprintf("CREATE TABLE %s%s ( %s );",
		mode, w.tableFullName, strings.Join(col, ", "))
	debug.Printf("SQL:%s", sql)
	_, err := w.Tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("%s: %s", err, sql)
	}
	return nil
}

func (w *Writer) createConstraints() []string {
	pk := tbln.SplitRow(toString(w.ExtraValue("primarykey")))
	nu := tbln.SplitRow(toString(w.ExtraValue("is_nullable")))
	if len(nu) != w.ColumnNum() {
		nu = nil
	}
	cs := make([]string, len(w.Names))
	for i := 0; i < len(w.Names); i++ {
		if (nu != nil) && nu[i] == "NO" {
			cs[i] += " NOT NULL"
		}
		if contains(pk, w.Names[i]) {
			cs[i] += " PRIMARY KEY"
		}
	}
	return cs
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func (w *Writer) prepare() error {
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
	// Construct SQL that does not generate an error
	// for each database when insert mode is OrIgnore.
	// MySQL
	// INSERT IGNORE INTO...
	// SQLlite3
	// INSERT OR IGNORE INTO...
	ignore := ""
	// PostgreSQL
	// INSERT INTO ... ON CONFLICT DO NOTHING
	onconf := ""
	if w.imode == OrIgnore {
		switch w.TDB.Name {
		case "mysql":
			ignore = "IGNORE "
		case "sqlite3":
			ignore = "OR IGNORE "
		case "postgres":
			onconf = "ON CONFLICT DO NOTHING"
		}
	}
	// #nosec G201
	insert := fmt.Sprintf(
		"INSERT %sINTO %s ( %s ) VALUES ( %s ) %s;",
		ignore,
		w.tableFullName, strings.Join(names, ", "), strings.Join(ph, ", "),
		onconf)
	debug.Printf("SQL:%s", insert)
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
