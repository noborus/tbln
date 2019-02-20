package db

import (
	"database/sql"
	"sync"
)

// Driver is the interface every database driver.
type Driver interface {
	GetPrimaryKey(db *sql.DB, tableName string) ([]string, error)
}

var drivers = make(map[string]Driver)
var driversMu sync.RWMutex

// Register is database driver register.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	drivers[name] = driver
}

// Get resturn database driver.
func Get(name string) Driver {
	return drivers[name]
}
