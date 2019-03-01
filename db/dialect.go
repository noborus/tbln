package db

import (
	"database/sql"
	"errors"
	"sync"
)

var (
	// ErrorNotSupport is database driver not supported
	ErrorNotSupport = errors.New("not supported")
)

// Driver is the interface every database driver.
type Driver struct {
	Style
	Constraint
}

// Style is SQLStyle struct
type Style struct {
	PlaceHolder string
	Quote       string
}

// Constraint is the interface database constraint.
type Constraint interface {
	GetPrimaryKey(db *sql.DB, TableName string) ([]string, error)
	GetColumnInfo(conn *sql.DB, TableName string) (map[string][]interface{}, error)
}

var drivers = make(map[string]Driver)
var driversMu sync.RWMutex

// Register is database driver register.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	drivers[name] = driver
}

// GetDriver resturn database driver.
func GetDriver(name string) Driver {
	return drivers[name]
}

// DefaultConstr is a driver that retrieves data using only database/SQL.
type DefaultConstr struct{}

// GetPrimaryKey returns the primary key as a slice.
func (c *DefaultConstr) GetPrimaryKey(conn *sql.DB, tableName string) ([]string, error) {
	return nil, ErrorNotSupport
}

// GetColumnInfo returns information of a table column as an array.
func (c *DefaultConstr) GetColumnInfo(conn *sql.DB, tableName string) (map[string][]interface{}, error) {
	return nil, ErrorNotSupport
}
