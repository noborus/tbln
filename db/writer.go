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
	ReplaceLN     bool
}

// NewWriter returns a new Writer that writes to database table.
func NewWriter(tdb *TDB, definition *tbln.Definition) (*Writer, error) {
	if !tdb.IsTx {
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
	wt := w.Types()
	for i, v := range row {
		// null
		if v == "" {
			r[i] = nil
			continue
		}
		if w.ReplaceLN {
			v = strings.ReplaceAll(v, "\\n", "\n")
		}
		r[i] = w.convertDBType(wt[i], v)
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
	if w.TableName() == "" {
		return fmt.Errorf("table name required")
	}
	if schema != "" {
		w.tableFullName = w.quoting(schema) + "." + w.quoting(w.TableName())
	} else {
		w.tableFullName = w.quoting(w.TableName())
	}
	err = w.WriteDefinition(cmode)
	if err != nil {
		return err
	}
	if cmode == CreateOnly {
		return nil
	}
	err = w.prepare(imode)
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

// WriteDefinition is create table and insert prepare.
func (w *Writer) WriteDefinition(cmode CreateMode) error {
	wn := w.Names()
	if wn == nil {
		if w.ColumnNum() == 0 {
			return fmt.Errorf("column num is 0")
		}
		wn = make([]string, w.ColumnNum())
		for i := 0; i < w.ColumnNum(); i++ {
			wn[i] = fmt.Sprintf("c%d", i+1)
		}
		err := w.SetNames(wn)
		if err != nil {
			return err
		}
	}
	wt := w.Types()
	if wt == nil {
		if w.ColumnNum() == 0 {
			return fmt.Errorf("column num is 0")
		}
		wt = make([]string, w.ColumnNum())
		for i := 0; i < w.ColumnNum(); i++ {
			wt[i] = "text"
		}
		err := w.SetTypes(wt)
		if err != nil {
			return err
		}
	}
	if cmode > NotCreate {
		err := w.createTable(cmode)
		if err != nil {
			return err
		}
	}
	return nil
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

func (w *Writer) createTable(cmode CreateMode) error {
	mode := ""
	if cmode == ReCreate {
		err := w.dropTable()
		if err != nil {
			return err
		}
	} else if cmode == IfNotExists {
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
		typeNames = w.Types()
	}
	wn := w.Names()
	col := make([]string, len(wn))
	for i := 0; i < len(wn); i++ {
		col[i] = w.quoting(wn[i]) + " " + typeNames[i] + constraints[i]
	}
	primaryKey := ""
	pk := tbln.SplitRow(toString(w.ExtraValue("primarykey")))
	if len(pk) > 0 {
		primaryKey = fmt.Sprintf(", PRIMARY KEY (%s)", strings.Join(pk, ","))
	}

	sql := fmt.Sprintf("CREATE TABLE %s%s ( %s %s);",
		mode, w.tableFullName, strings.Join(col, ", "), primaryKey)
	debug.Printf("SQL:%s", sql)
	_, err := w.Tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("%s: %s", err, sql)
	}
	return nil
}

func (w *Writer) createConstraints() []string {
	nu := tbln.SplitRow(toString(w.ExtraValue("is_nullable")))
	if len(nu) != w.ColumnNum() {
		nu = nil
	}
	wn := w.Names()
	cs := make([]string, len(wn))
	for i := 0; i < len(wn); i++ {
		if (nu != nil) && nu[i] == "NO" {
			cs[i] += " NOT NULL"
		}
	}
	return cs
}

func (w *Writer) prepare(imode InsertMode) error {
	var err error
	wn := w.Names()
	names := make([]string, len(wn))
	ph := make([]string, len(wn))
	for i := 0; i < len(wn); i++ {
		names[i] = w.quoting(wn[i])
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
	if imode == OrIgnore {
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
	case "bool":
		if w.TDB.Name == "mysql" {
			if value == "false" || value == "f" || value == "0" {
				return 0
			}
			return 1
		}
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
