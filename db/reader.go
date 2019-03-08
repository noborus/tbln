package db

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/noborus/tbln"
)

// Reader reads records from database table.
type Reader struct {
	*tbln.Definition
	*TDB
	rows     *sql.Rows
	scanArgs []interface{}
	values   []interface{}
}

// ReadTable returns a new Reader from table name.
func (tdb *TDB) ReadTable(schema string, tableName string, pkey []string) (*Reader, error) {
	tr := &Reader{
		Definition: tbln.NewDefinition(),
		TDB:        tdb,
	}
	tr.SetTableName(tableName)
	var err error
	if schema == "" {
		schema, err = tr.GetSchema(tdb.DB)
		if err != nil {
			return nil, err
		}
	}
	// Constraint
	info, err := tr.GetColumnInfo(tdb.DB, schema, tableName)
	if err != nil {
		if err != ErrorNotSupport {
			return nil, err
		}
	}
	tr.setTableInfo(info)
	// Primary key
	pk, err := tr.GetPrimaryKey(tr.TDB.DB, schema, tableName)
	if err != nil && err != ErrorNotSupport {
		return nil, err
	} else if len(pk) > 0 {
		tr.Extras["primarykey"] = tbln.NewExtra(tbln.JoinRow(pk), false)
	}
	if len(pkey) == 0 && len(pk) > 0 {
		pkey = pk
	}
	var orderby string
	if len(pkey) > 0 {
		orderby = strings.Join(pkey, ", ")
	} else {
		orderby = "1"
	}
	table := tdb.quoting(tableName)
	if schema != "" {
		table = tdb.quoting(schema) + "." + tdb.quoting(tableName)
	}
	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY %s", table, orderby)
	err = tr.query(sql)
	if err != nil {
		return nil, fmt.Errorf("%s: [%s]", err, sql)
	}
	return tr, nil
}

// ReadQuery returns a new Reader from SQL query.
func (tdb *TDB) ReadQuery(sql string, args ...interface{}) (*Reader, error) {
	tr := &Reader{
		Definition: tbln.NewDefinition(),
		TDB:        tdb,
	}
	err := tr.query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: (%s)", err, sql)
	}
	return tr, nil
}

// ReadRow reads one record (a slice of fields) from tr.
func (tr *Reader) ReadRow() ([]string, error) {
	if !tr.rows.Next() {
		err := tr.rows.Err()
		return nil, err
	}
	err := tr.rows.Scan(tr.scanArgs...)
	if err != nil {
		return nil, err
	}
	rec := make([]string, len(tr.values))
	for i, col := range tr.values {
		rec[i] = toString(col)
	}
	return rec, nil
}

// ReadTableAll reads all the remaining records from tableName.
func ReadTableAll(tdb *TDB, schema string, tableName string) (*tbln.Tbln, error) {
	tr, err := tdb.ReadTable(schema, tableName, nil)
	if err != nil {
		return nil, err
	}
	return tr.readRowsAll()
}

// ReadQueryAll reads all the remaining records from SQL query.
func ReadQueryAll(tdb *TDB, query string, args ...interface{}) (*tbln.Tbln, error) {
	tr, err := tdb.ReadQuery(query, args...)
	if err != nil {
		return nil, err
	}
	return tr.readRowsAll()
}

func (tr *Reader) readRowsAll() (*tbln.Tbln, error) {
	var err error
	defer func() {
		cerr := tr.rows.Close()
		if err == nil {
			err = cerr
		}
	}()

	at := &tbln.Tbln{}
	at.Definition = tr.Definition
	at.Rows = make([][]string, 0)
	for {
		var rec []string
		rec, err = tr.ReadRow()
		if err != nil {
			return at, err
		}
		if rec == nil {
			break
		}
		at.RowNum++
		at.Rows = append(at.Rows, rec)
	}
	return at, err
}

func (tr *Reader) setTableInfo(constraints map[string][]interface{}) {
	for k, v := range constraints {
		col := make([]string, len(v))
		visible := false
		for i, c := range v {
			col[i] = toString(c)
			if col[i] != "" {
				visible = true
			}
		}
		if visible {
			tr.Extras[strings.ToLower(k)] = tbln.NewExtra(tbln.JoinRow(col), false)
		}
	}
}

func (tr *Reader) query(query string, args ...interface{}) error {
	rows, err := tr.DB.Query(query, args...)
	if err != nil {
		return err
	}

	err = tr.setRowInfo(rows)
	if err != nil {
		return err
	}
	tr.values = make([]interface{}, tr.ColumnNum())
	tr.scanArgs = make([]interface{}, tr.ColumnNum())
	for i := range tr.values {
		tr.scanArgs[i] = &tr.values[i]
	}
	tr.rows = rows
	return nil
}

func (tr *Reader) setRowInfo(rows *sql.Rows) error {
	var err error
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	err = tr.SetNames(columns)
	if err != nil {
		return err
	}
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
	err = tr.SetTypes(types)
	if err != nil {
		return err
	}
	if _, ok := tr.Extras[tr.Name+"_type"]; !ok {
		tr.Extras[tr.Name+"_type"] = tbln.NewExtra(tbln.JoinRow(dbtypes), false)
	}
	return nil
}

func toString(v interface{}) string {
	var str string
	switch t := v.(type) {
	case nil:
		str = ""
	case time.Time:
		str = t.Format(time.RFC3339)
	case []byte:
		if ok := utf8.Valid(t); ok {
			str = string(t)
		} else {
			str = `\x` + hex.EncodeToString(t)
		}
	default:
		str = fmt.Sprint(v)
		str = strings.ReplaceAll(str, "\n", "\\n")
	}
	return str
}

func convertType(dbtype string) string {
	switch strings.ToLower(dbtype) {
	case "smallint", "integer", "int", "int2", "int4", "smallserial", "serial":
		return "int"
	case "bigint", "int8", "bigserial":
		return "bigint"
	case "float", "decimal", "numeric", "real", "double precision":
		return "numeric"
	case "bool":
		return "bool"
	case "timestamp", "timestamptz", "date", "time":
		return "timestamp"
	case "string", "text", "char", "varchar":
		return "text"
	default:
		return "text"
	}
}
