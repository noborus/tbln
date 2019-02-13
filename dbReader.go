package tbln

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// DBReader is DB read struct.
type DBReader struct {
	Definition
	*DBD
	db       *sql.DB
	tx       *sql.Tx
	rows     *sql.Rows
	scanArgs []interface{}
	values   []interface{}
}

// NewDBReader is creates a structure for reading from the DB table.
func NewDBReader(dbd *DBD, tableName string) (*DBReader, error) {
	tr := &DBReader{
		Definition: NewDefinition(tableName),
		DBD:        dbd,
		db:         dbd.DB,
	}
	err := tr.preparation()
	if err != nil {
		return nil, err
	}
	return tr, nil
}

// ReadRow is return one row.
func (tr *DBReader) ReadRow() ([]string, error) {
	if !tr.rows.Next() {
		err := tr.rows.Err()
		tr.rows.Close()
		return nil, err
	}
	err := tr.rows.Scan(tr.scanArgs...)
	if err != nil {
		return nil, err
	}
	rec := make([]string, len(tr.values))
	for i, col := range tr.values {
		rec[i] = valString(col)
	}
	return rec, nil
}

func valString(v interface{}) string {
	var str string
	b, ok := v.([]byte)
	if ok {
		str = string(b)
	} else {
		if v == nil {
			str = ""
		} else {
			str = fmt.Sprint(v)
		}
	}
	return str
}

func (tr *DBReader) preparation() error {
	rows, err := tr.Query(`SELECT * FROM ` + tr.tableName)
	if err != nil {
		return err
	}
	err = tr.setInfo(rows)
	if err != nil {
		return err
	}
	tr.values = make([]interface{}, tr.columnNum)
	tr.scanArgs = make([]interface{}, tr.columnNum)
	for i := range tr.values {
		tr.scanArgs[i] = &tr.values[i]
	}
	tr.rows = rows
	return nil
}

// Query is sql.Query wrapper.
func (tr *DBReader) Query(query string) (*sql.Rows, error) {
	return tr.db.Query(query)
}

func (tr *DBReader) setInfo(rows *sql.Rows) error {
	var err error
	tr.SetTableName(tr.tableName)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	tr.SetNames(columns)
	columntype, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	dbtypes := make([]string, len(columns))
	types := make([]string, len(columns))
	for i, ct := range columntype {
		dbtypes[i] = ct.DatabaseTypeName()
		types[i] = convertType(ct.DatabaseTypeName())
	}
	tr.SetTypes(types)
	// Database type
	t := "| " + strings.Join(dbtypes, " | ") + " |"
	tr.Ext[tr.Name+"_type"] = Extra{value: t, hashTarget: true}
	return nil
}

func convertType(dbtype string) string {
	switch strings.ToLower(dbtype) {
	case "smallint", "integer", "int", "int4", "smallserial", "serial":
		return "int"
	case "bigint", "int8", "bigserial":
		return "bigint"
	case "float", "decimal", "numeric", "real", "double precision":
		return "numeric"
	case "bool":
		return "bool"
	case "string", "text", "char", "varchar":
		return "text"
	case "timestamp", "timestamptz", "date", "time":
		return "timestamp"
	default:
		return "text"
	}
}

// ReadTable reads all rows in the table.
func ReadTable(db *DBD, tableName string) (*Table, error) {
	r, err := NewDBReader(db, tableName)
	if err != nil {
		return nil, err
	}
	at := &Table{}
	at.Definition = r.Definition
	at.Rows = make([][]string, 0)
	for {
		rec, err := r.ReadRow()
		if err != nil {
			log.Println(err)
			break
		}
		if rec == nil {
			break
		}
		at.RowNum++
		at.Rows = append(at.Rows, rec)
	}
	return at, err
}
