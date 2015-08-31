package forms

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotANumber = errors.New("not a number")
)

type LoaderFunc func(string) (interface{}, error)

var StringLoader = LoaderFunc(func(rawValue string) (interface{}, error) {
	return strings.TrimSpace(rawValue), nil
})

var IntLoader = LoaderFunc(func(rawValue string) (interface{}, error) {
	val, err := strconv.ParseInt(rawValue, 0, 0)
	if err != nil {
		return nil, ErrNotANumber
	}
	return int(val), nil
})

var TimeLoader = NewTimeLoader(time.RFC3339)

func NewTimeLoader(layout string) LoaderFunc {
	return LoaderFunc(func(rawValue string) (interface{}, error) {
		val, err := time.Parse(layout, rawValue)
		if err != nil {
			return nil, err
		}
		return val, nil
	})
}
