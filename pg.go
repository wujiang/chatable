package chatable

import (
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"regexp"
	"strings"
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
	case time.Time:
		nt.Time, nt.Valid = src.(time.Time)
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

// StringSlice is one dimension text array in Postgresql
// https://gist.github.com/adharris/4163702#gistcomment-1356268
type StringSlice []string

var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Scan convert to a slice of strings
// http://www.postgresql.org/docs/9.1/static/arrays.html#ARRAYS-IO
func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	str := string(asBytes)

	// change quote escapes for csv parser
	str = quoteEscapeRegex.ReplaceAllString(str, `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)
	// remove braces
	str = str[1 : len(str)-1]
	csvReader := csv.NewReader(strings.NewReader(str))

	slice, err := csvReader.Read()

	if err != nil {
		return err
	}

	(*s) = StringSlice(slice)

	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	for i, elem := range s {
		s[i] = `"` + strings.Replace(strings.Replace(elem, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(s, ",") + "}", nil
}
