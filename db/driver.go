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

	Constraint
}

// Constraint is the interface database constraint.
type Constraint interface {
	GetPrimaryKey(db *sql.DB, TableName string) ([]string, error)
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

// Default is a driver that retrieves data using only database/SQL.
type Default struct{}

// GetPrimaryKey returns the primary key as a slice.
func (d *Default) GetPrimaryKey(db *sql.DB, TableName string) ([]string, error) {
	return nil, fmt.Errorf("this database is not supported")
}

// PlaceHolder returns the placeholer string.
func (d *Default) PlaceHolder() string {
	return "?"
}

// Quote returns the quote string.
func (d *Default) Quote() string {
	return `"`
}
