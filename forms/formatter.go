package forms

import "strings"

type FormatterFunc func(string) string

var CapFormatter = FormatterFunc(func(rawValue string) string {
	return strings.Title(strings.ToLower(rawValue))
})
