package tbln

import (
	"database/sql"
	"log"
)

// DBExport is output in tbln format from DB.
func (tbln *Tbln) DBExport(db *sql.DB, table string) error {
	sql := `SELECT * FROM ` + table
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	rows, err := tx.Query(sql)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	tbln.AddExtra("TableName: " + table)
	tbln.SetNames(columns)
	columntype, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	types := make([]string, len(columns))
	for i, ct := range columntype {
		types[i] = ct.DatabaseTypeName()
	}
	tbln.SetTypes(types)
	values := make([]string, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		tbln.AddRow(values)
	}
	err = rows.Err()
	tx.Commit()
	return nil
}
