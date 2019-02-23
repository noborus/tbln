package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/noborus/tbln"
)

// Reader is DB read struct.
type Reader struct {
	tbln.Definition
	query string
	*TDB
	rows     *sql.Rows
	scanArgs []interface{}
	values   []interface{}
}

// ReadTable is reates a structure for reading from DB table.
func (TDB *TDB) ReadTable(TableName string, pkey []string) (*Reader, error) {
	tr := &Reader{
		Definition: tbln.NewDefinition(),
		TDB:        TDB,
	}
	tr.SetTableName(TableName)
	// Constraint
	columns, err := tr.GetColumnInfo(tr.TDB.DB, tr.TableName)
	if err != nil {
		if err != ErrorNotSupport {
			return nil, err
		}
	} else if columns != nil {
		tr.constraint(columns)
	}
	// Primary key
	pk, err := tr.GetPrimaryKey(tr.TDB.DB, tr.TableName)
	if err != nil && err != ErrorNotSupport {
		return nil, err
	} else if len(pk) > 0 {
		tr.Ext["primarykey"] = tbln.NewExtra(tbln.JoinRow(pk))
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
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s", TDB.quoting(TableName), orderby)
	err = tr.peparation(query)
	if err != nil {
		return nil, fmt.Errorf("%s: [%s]", err, query)
	}
	return tr, nil
}

// ReadQuery is reates a structure for reading from query.
func (TDB *TDB) ReadQuery(query string, args ...interface{}) (*Reader, error) {
	tr := &Reader{
		Definition: tbln.NewDefinition(),
		TDB:        TDB,
	}
	err := tr.peparation(query, args...)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

// ReadRow is return one row.
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
		rec[i] = valString(col)
	}
	return rec, nil
}

// ReadTableAll reads all rows in the table.
func ReadTableAll(TDB *TDB, TableName string) (*tbln.Tbln, error) {
	r, err := TDB.ReadTable(TableName, nil)
	if err != nil {
		return nil, err
	}
	return readRowsAll(r)
}

// ReadQueryAll reads all rows in the table.
func ReadQueryAll(TDB *TDB, query string, args ...interface{}) (*tbln.Tbln, error) {
	r, err := TDB.ReadQuery(query, args...)
	if err != nil {
		return nil, err
	}
	return readRowsAll(r)
}

func readRowsAll(rd *Reader) (at *tbln.Tbln, err error) {
	at = &tbln.Tbln{}
	at.Definition = rd.Definition
	at.Rows = make([][]string, 0)
	defer func() {
		cerr := rd.rows.Close()
		if err == nil {
			err = cerr
		}
	}()
	for {
		var rec []string
		rec, err = rd.ReadRow()
		if err != nil {
			return at, err
		}
		if rec == nil {
			break
		}
		at.RowNum++
		at.Rows = append(at.Rows, rec)
	}
	return
}

func (tr *Reader) constraint(columns map[string][]interface{}) (map[string]string, error) {
	for k, v := range columns {
		col := make([]string, len(v))
		visible := false
		for i, c := range v {
			col[i] = valString(c)
			if col[i] != "" {
				visible = true
			}
		}
		if visible {
			tr.Ext[strings.ToLower(k)] = tbln.NewExtra(tbln.JoinRow(col))
		}
	}
	return nil, nil
}

// preparation is read preparation.
func (tr *Reader) peparation(query string, args ...interface{}) error {
	tr.query = query
	rows, err := tr.DB.Query(query, args...)
	if err != nil {
		return err
	}
	err = tr.setExtra(rows)
	if err != nil {
		return err
	}
	tr.values = make([]interface{}, tr.ColumnNum)
	tr.scanArgs = make([]interface{}, tr.ColumnNum)
	for i := range tr.values {
		tr.scanArgs[i] = &tr.values[i]
	}
	tr.rows = rows
	return nil
}

func (tr *Reader) setExtra(rows *sql.Rows) error {
	var err error
	tr.SetTableName(tr.TableName)
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
	// Database type
	if _, ok := tr.Ext[tr.Name+"_type"]; !ok {
		tr.Ext[tr.Name+"_type"] = tbln.NewExtra(tbln.JoinRow(dbtypes))
	}
	return nil
}

func valString(v interface{}) string {
	var str string
	b, ok := v.([]byte)
	if ok {
		str = string(b)
	} else {
		switch t := v.(type) {
		case nil:
			str = ""
		case time.Time:
			str = t.Format(time.RFC3339)
		default:
			str = fmt.Sprint(v)
		}

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
	case "string", "text", "char", "varchar":
		return "text"
	case "timestamp", "timestamptz", "date", "time":
		return "timestamp"
	default:
		return "text"
	}
}
