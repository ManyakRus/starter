package carbon

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Scan an interface used by Scan in package database/sql for Scanning value from database to local golang variable.
func (c *Carbon) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*c = Carbon{time: value, loc: time.Local}
		return nil
	}
	return fmt.Errorf("can not convert %v to carbon", v)
}

// Value the interface providing the Value method for package database/sql/driver.
func (c Carbon) Value() (driver.Value, error) {
	if c.IsZero() {
		return nil, nil
	}
	return c.time, nil
}
