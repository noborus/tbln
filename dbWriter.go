package tbln

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
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
func NewDBWriter(dbd *DBD, d Definition, create bool) *DBWriter {
	return &DBWriter{
		Definition: d,
		DBD:        dbd,
		db:         dbd.DB,
		Create:     create,
	}
}

// WriteDefinition is write table definition.
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
		tw.quoting(tw.name), strings.Join(col, ", "))
	_, err := tw.db.Exec(sql)
	fmt.Println(sql)
	if err != nil {
		return err
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
		tw.quoting(tw.name), strings.Join(names, ", "), strings.Join(ph, ", "))
	tw.stmt, err = tw.db.Prepare(insert)
	fmt.Println(insert)
	return err
}

func (tw *DBWriter) quoting(name string) string {
	r := regexp.MustCompile(`[^a-z0-9_]+`)
	if r.MatchString(name) {
		return tw.Quote + name + tw.Quote
	}
	return name
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

// WriteTable writes all rows to the table.
func WriteTable(db *DBD, table *Table, create bool) error {
	w := NewDBWriter(db, table.Definition, create)
	err := w.WriteDefinition()
	if err != nil {
		return err
	}
	for _, row := range table.Rows {
		err = w.WriteRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}
