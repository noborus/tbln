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
	rows      *sql.Rows
	scanArgs  []interface{}
	values    []interface{}
	mySQLType []string
}

// RangePrimaryKey is range of primary key.
type RangePrimaryKey struct {
	name string
	min  interface{}
	max  interface{}
}

// NewRangePrimaryKey return RangePrimaryKey
func NewRangePrimaryKey(name string, min interface{}, max interface{}) RangePrimaryKey {
	return RangePrimaryKey{name: name, min: min, max: max}
}

// ReadTable returns a new Reader from table name.
func (tdb *TDB) ReadTable(schema string, tableName string, rps []RangePrimaryKey) (*Reader, error) {
	if tableName == "" {
		return nil, fmt.Errorf("require table name")
	}
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
	tr.mySQLType = tbln.SplitRow(toString(tr.ExtraValue("mysql_columntype")))
	// Primary key
	pkey, err := tr.GetPrimaryKey(tr.TDB.DB, schema, tableName)
	if err != nil && err != ErrorNotSupport {
		return nil, err
	}
	if len(pkey) > 0 {
		tr.Extras["primarykey"] = tbln.NewExtra(tbln.JoinRow(pkey), false)
	}
	table := tr.fullTableName(schema, tableName)
	conds, args := tr.conditions(pkey, rps)
	order := tr.orderby(pkey)

	// #nosec G201
	sql := fmt.Sprintf("SELECT * FROM %s %s %s", table, conds, order)
	debug.Printf("SQL:%s:%s", sql, args)
	err = tr.query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: [%s]", err, sql)
	}
	return tr, nil
}

func (tr *Reader) conditions(pkey []string, rps []RangePrimaryKey) (string, []interface{}) {
	if len(rps) == 0 {
		return "", nil
	}
	if len(pkey) != len(rps) {
		debug.Printf("primary key miss match!")
		return "", nil
	}

	args := make([]interface{}, 0, len(rps)*2)
	var cs []string
	for i, rp := range rps {
		if pkey[i] != rp.name {
			debug.Printf("primary key miss match! %s:%s", pkey[i], rp.name)
			return "", nil
		}
		if tr.Style.PlaceHolder == "$" {
			cd := fmt.Sprintf("(%s >= $%d AND %s <= $%d)", tr.quoting(rp.name), i*2+1, tr.quoting(rp.name), i*2+2)
			cs = append(cs, cd)
		} else {
			cd := fmt.Sprintf("(%s >= ? AND %s <= ?)", tr.quoting(rp.name), tr.quoting(rp.name))
			cs = append(cs, cd)
		}
		args = append(args, rp.min)
		args = append(args, rp.max)
	}
	// #nosec G202
	conds := " WHERE " + strings.Join(cs, " AND ")
	return conds, args
}

func (tr *Reader) orderby(pkey []string) string {
	if len(pkey) == 0 {
		return "ORDER BY 1"
	}
	pk := make([]string, len(pkey))
	for i, p := range pkey {
		pk[i] = tr.quoting(p)
	}
	return "ORDER BY " + strings.Join(pk, ", ")
}

// ReadQuery returns a new Reader from SQL query.
func (tdb *TDB) ReadQuery(sql string, args ...interface{}) (*Reader, error) {
	tr := &Reader{
		Definition: tbln.NewDefinition(),
		TDB:        tdb,
	}
	debug.Printf("SQL:%s:%s", sql, args)
	err := tr.query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: (%s)", err, sql)
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
	// MySQL converts to string by columnType,
	// because MySQL drivers are returned in all []byte.
	if tr.TDB.Name == "mysql" {
		var mTypes []string
		if len(tr.mySQLType) == len(tr.values) {
			mTypes = tr.mySQLType
		} else {
			mTypes = tr.Types()
		}
		for i, col := range tr.values {
			rec[i] = mySQLtoString(mTypes[i], col)
		}
		return rec, nil
	}
	for i, col := range tr.values {
		rec[i] = toString(col)
	}
	return rec, nil
}

// Close reader.
func (tr *Reader) Close() error {
	return tr.rows.Close()
}

// GetTableInfo returns only the information of table.
func GetTableInfo(tdb *TDB, schema string, tableName string) (*tbln.TBLN, error) {
	tr, err := tdb.ReadTable(schema, tableName, nil)
	if err != nil {
		return nil, err
	}
	at := &tbln.TBLN{}
	at.Definition = tr.Definition
	return at, nil
}

// ReadTableAll reads all the remaining records from tableName.
func ReadTableAll(tdb *TDB, schema string, tableName string) (*tbln.TBLN, error) {
	tr, err := tdb.ReadTable(schema, tableName, nil)
	if err != nil {
		return nil, err
	}
	return tr.readRowsAll()
}

// ReadQueryAll reads all the remaining records from SQL query.
func ReadQueryAll(tdb *TDB, query string, args ...interface{}) (*tbln.TBLN, error) {
	tr, err := tdb.ReadQuery(query, args...)
	if err != nil {
		return nil, err
	}
	return tr.readRowsAll()
}

func (tr *Reader) readRowsAll() (*tbln.TBLN, error) {
	var err error
	defer func() {
		closeErr := tr.Close()
		if err == nil {
			err = closeErr
		}
	}()

	at := &tbln.TBLN{}
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
	dbTypes := make([]string, len(columns))
	types := make([]string, len(columns))
	for i, ct := range columntype {
		dbTypes[i] = ct.DatabaseTypeName()
		types[i] = convertType(ct.DatabaseTypeName())
	}

	// Convert MySQL tinyint(1) to bool type
	if tr.TDB.Name == "mysql" {
		tr.mySQLType = tbln.SplitRow(toString(tr.ExtraValue("mysql_columntype")))
		for i, t := range tr.mySQLType {
			if t == "tinyint(1)" {
				types[i] = "bool"
			}
		}
	}

	err = tr.SetTypes(types)
	if err != nil {
		return err
	}
	if _, ok := tr.Extras[tr.Name+"_type"]; !ok {
		tr.Extras[tr.Name+"_type"] = tbln.NewExtra(tbln.JoinRow(dbTypes), false)
	}
	return nil
}

func mySQLtoString(dbType string, v interface{}) string {
	const layout = "2006-01-02 15:04:05"
	var str string

	if b, ok := v.([]byte); ok {
		if ok := utf8.Valid(b); !ok {
			str = `\x` + hex.EncodeToString(b)
			return str
		}
		str = string(b)
	} else {
		str = fmt.Sprint(v)
	}
	switch dbType {
	case "tinyint(1)":
		if str == "1" {
			str = "true"
		} else if str == "0" {
			str = "false"
		}
	case "timestamp":
		t, err := time.Parse(layout, str)
		if err != nil {
			return ""
		}
		str = t.Format(time.RFC3339)
	}
	return str
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

func convertType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "tinyint", "smallint", "integer", "int", "int2", "int4", "smallserial", "serial":
		return "int"
	case "bigint", "int8", "bigserial":
		return "bigint"
	case "float4", "float", "real":
		return "double precision"
	case "double", "double precision", "float8":
		return "double precision"
	case "decimal", "numeric":
		return "numeric"
	case "bool":
		return "bool"
	case "timestamp", "timestamptz", "date", "time":
		return "timestamp"
	case "string", "text", "char", "varchar", "character", "bpchar":
		return "text"
	default:
		debug.Printf("unsupported type:%s", dbType)
		return "text"
	}
}
