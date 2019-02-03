package tbln

import (
	"database/sql"
	"log"
)

// DBReader is DB read struct.
type DBReader struct {
	Table
	db   *sql.DB
	tx   *sql.Tx
	rows *sql.Rows
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
		err := tr.readInfo()
		if err != nil {
			return nil, err
		}
	}
	values := make([]string, tr.columnNum)
	scanArgs := make([]interface{}, tr.columnNum)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for tr.rows.Next() {
		err := tr.rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		return values, nil
	}
	err := tr.rows.Err()
	return nil, err
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
		types[i] = ct.DatabaseTypeName()
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
