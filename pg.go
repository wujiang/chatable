package asapp

import (
	"database/sql/driver"
	"errors"
	"time"
)

// gorp use []bytes as src.

// NullTime represents a time.Time that may be null. NullTime implements the
// sql.Scanner interface so it can be used as a scan destination, similar to
// sql.NullString.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		dt, err := time.Parse(time.RFC3339Nano, string(s))
		nt.Time, nt.Valid = dt, err == nil
	case nil:
		nt.Valid = false
	default:
		return error(errors.New("Scan source was not []bytes or nil"))
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
