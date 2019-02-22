package db

import (
	"database/sql"
	"fmt"
	"sync"
)

// Driver is the interface every database driver.
type Driver interface {
	PlaceHolder() string
	Quote() string
}

// Constraint is the interface database constraint.
type Constraint interface {
	GetPrimaryKey(db *sql.DB, TableName string) ([]string, error)
	GetColumnInfo(conn *sql.DB, TableName string) (map[string][]interface{}, error)
}

var drivers = make(map[string]Driver)
var constraints = make(map[string]Constraint)
var driversMu sync.RWMutex

// Register is database driver register.
func Register(name string, driver Driver, constraint Constraint) {
	driversMu.Lock()
	defer driversMu.Unlock()

	drivers[name] = driver
	if constraint != nil {
		constraints[name] = constraint
	}
}

// GetDriver resturn database driver.
func GetDriver(name string) (Driver, Constraint) {
	return drivers[name], constraints[name]
}

// Default is a driver that retrieves data using only database/SQL.
type Default struct{}

// PlaceHolder returns the placeholer string.
func (d *Default) PlaceHolder() string {
	return "?"
}

// Quote returns the quote string.
func (d *Default) Quote() string {
	return `"`
}

// DefaultConstr is a driver that retrieves data using only database/SQL.
type DefaultConstr struct{}

// GetPrimaryKey returns the primary key as a slice.
func (d *DefaultConstr) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	return nil, fmt.Errorf("not supported")
}

// GetColumnInfo returns information of a table column as an array.
func (d *DefaultConstr) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	return nil, fmt.Errorf("not supported")
}
