package db

import "database/sql"

// Constraint is get the Constraint for each database
type Constraint interface {
	GetPrimaryKey(db *sql.DB, tableName string) ([]string, error)
}

var Driver string
var DB map[string]Constraint = make(map[string]Constraint)
