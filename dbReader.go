package tbln

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// DBReader is DB read struct.
type DBReader struct {
	Definition
	query string
	*DBD
	rows     *sql.Rows
	scanArgs []interface{}
	values   []interface{}
}

// ReadTable is reates a structure for reading from DB table.
func (dbd *DBD) ReadTable(tableName string, pkey []string) (*DBReader, error) {
	tr := &DBReader{
		Definition: NewDefinition(),
		DBD:        dbd,
	}
	tr.SetTableName(tableName)
	var orderby string
	if len(pkey) > 0 {
		orderby = strings.Join(pkey, ", ")
	} else {
		orderby = "1"
	}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s", dbd.quoting(tableName), orderby)
	err := tr.peparation(query)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

// ReadQuery is reates a structure for reading from query.
func (dbd *DBD) ReadQuery(query string, args ...interface{}) (*DBReader, error) {
	tr := &DBReader{
		Definition: NewDefinition(),
		DBD:        dbd,
	}
	err := tr.peparation(query, args...)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

// ReadRow is return one row.
func (tr *DBReader) ReadRow() ([]string, error) {
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

// preparation is read preparation.
func (tr *DBReader) peparation(query string, args ...interface{}) error {
	tr.query = query
	rows, err := tr.DB.Query(query, args...)
	if err != nil {
		return err
	}
	err = tr.setExtra(rows)
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

func (tr *DBReader) setExtra(rows *sql.Rows) error {
	var err error
	tr.SetTableName(tr.tableName)
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
	tr.Ext[tr.Name+"_type"] = NewExtra(stringRow(dbtypes))
	// Primary key
	pk, err := tr.DBD.GetPrimaryKey(tr.DBD.DB, tr.tableName)
	if len(pk) > 0 && err == nil {
		tr.Ext["Primarykey"] = NewExtra(stringRow(pk))
	}

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
	r, err := db.ReadTable(tableName, nil)
	if err != nil {
		return nil, err
	}
	return readRowsAll(db, r)
}

// ReadQueryAll reads all rows in the table.
func ReadQueryAll(db *DBD, query string, args ...interface{}) (*Tbln, error) {
	r, err := db.ReadQuery(query, args...)
	if err != nil {
		return nil, err
	}
	return readRowsAll(db, r)
}

func readRowsAll(db *DBD, rd *DBReader) (at *Tbln, err error) {
	at = &Tbln{}
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
