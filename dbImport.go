package tbln

import (
	"database/sql"
	"fmt"
	"strings"
)

// DBImport is input in tbln format to DB.
func (tbln *Tbln) DBImport(db *sql.DB, table string) error {
	err := tbln.createTable(db, table)
	if err != nil {
		return err
	}
	presql := "INSERT INTO " + table + " ("
	presql += strings.Join(tbln.Names, ", ")
	presql += ") VALUES ("
	for _, row := range tbln.Rows {
		sql := presql + strings.Join(row, ", ") + ");"
		fmt.Println(sql)
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tbln *Tbln) createTable(db *sql.DB, table string) error {
	sql := "CREATE TABLE " + table + " ("
	col := make([]string, len(tbln.Names))
	for i := 0; i < len(tbln.Names); i++ {
		col[i] = tbln.Names[i] + " " + tbln.Types[i]
	}
	sql += strings.Join(col, ", ")
	sql += ");"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	fmt.Println(sql)
	return nil
}
