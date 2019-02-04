package tbln

import (
	"database/sql"
	"log"
	"strings"
)

// DBReader is DB read struct.
type DBReader struct {
	Table
	db       *sql.DB
	tx       *sql.Tx
	rows     *sql.Rows
	scanArgs []interface{}
	values   []string
}

// NewDBReader is creates a structure for reading from the DB table.
func NewDBReader(db *sql.DB, tableName string) *DBReader {
	return &DBReader{
		Table: Table{Ext: make(map[string]string), name: tableName},
		db:    db,
	}
}

// ReadRow is return one row.
func (tr *DBReader) ReadRow() ([]string, error) {
	if tr.rows == nil {
		err := tr.preparation()
		if err != nil {
			return nil, err
		}
	}
	if !tr.rows.Next() {
		return nil, tr.rows.Err()
	}
	err := tr.rows.Scan(tr.scanArgs...)
	if err != nil {
		log.Fatal(err)
	}
	return tr.values, nil
}

func (tr *DBReader) preparation() error {
	err := tr.readInfo()
	if err != nil {
		return err
	}
	tr.values = make([]string, tr.columnNum)
	tr.scanArgs = make([]interface{}, tr.columnNum)
	for i := range tr.values {
		tr.scanArgs[i] = &tr.values[i]
	}
	return nil
}

func (tr *DBReader) readInfo() error {
	var err error
	tr.Ext["TableName"] = tr.name
	err = tr.begin()
	if err != nil {
		return err
	}
	rows, err := tr.tx.Query(`SELECT * FROM ` + tr.name)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	tr.setNames(columns)
	columntype, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	types := make([]string, len(columns))
	for i, ct := range columntype {
		types[i] = convertType(ct.DatabaseTypeName())
	}
	tr.setTypes(types)
	tr.rows = rows
	return nil
}

// Close is close the cursor.
func (tr *DBReader) Close() error {
	if tr.rows != nil {
		err := tr.rows.Close()
		tr.rows = nil
		return err
	}
	return tr.commit()

}

func (tr *DBReader) begin() error {
	var err error
	tr.tx, err = tr.db.Begin()
	if err != nil {
		return err
	}
	return nil
}

func (tr *DBReader) commit() error {
	return tr.tx.Commit()
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
	default:
		return "text"
	}
}
