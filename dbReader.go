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
	table string
	query string
	*DBD
	db       *sql.DB
	tx       *sql.Tx
	rows     *sql.Rows
	scanArgs []interface{}
	values   []interface{}
}

// ReadTable is reates a structure for reading from DB table.
func (dbd *DBD) ReadTable(tableName string) (*DBReader, error) {
	r := &DBReader{
		Definition: NewDefinition(),
		DBD:        dbd,
		db:         dbd.DB,
	}
	r.SetTableName(tableName)
	query := `SELECT * FROM ` + tableName + ` ORDER BY 1`
	err := r.peparation(query)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ReadQuery is reates a structure for reading from query.
func (dbd *DBD) ReadQuery(query string, args ...interface{}) (*DBReader, error) {
	r := &DBReader{
		Definition: NewDefinition(),
		DBD:        dbd,
		db:         dbd.DB,
	}
	err := r.peparation(query, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
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

// preparation is read preparation.
func (tr *DBReader) peparation(query string, args ...interface{}) error {
	rows, err := tr.db.Query(query, args...)
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
	tr.Ext[tr.Name+"_type"] = Extra{value: t, hashTarget: false}
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

// ReadTableAll reads all rows in the table.
func ReadTableAll(db *DBD, tableName string) (*Tbln, error) {
	r, err := db.ReadTable(tableName)
	if err != nil {
		return nil, err
	}
	return readTableAll(db, r)
}

// ReadQueryAll reads all rows in the table.
func ReadQueryAll(db *DBD, query string, args ...interface{}) (*Tbln, error) {
	r, err := db.ReadQuery(query, args...)
	if err != nil {
		return nil, err
	}
	return readTableAll(db, r)
}

func readTableAll(db *DBD, rd *DBReader) (*Tbln, error) {
	at := &Tbln{}
	at.Definition = rd.Definition
	at.Rows = make([][]string, 0)
	var err error
	for {
		rec, err := rd.ReadRow()
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
