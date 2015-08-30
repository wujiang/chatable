package forms

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithFormatters(t *testing.T) {
	fb := FieldBuilder{
		formatters: []FormatterFunc{},
		validators: []ValidatorFunc{},
	}
	f := FormatterFunc(func(s string) string { return s })
	assert.Equal(t, fb.formatters, []FormatterFunc{})
	fb.WithFormatters(f)
	// reflect.DeepEqual on the same function returns false
	// for example:
	// var tf = func(s string) string { return s }
	// assert.Equal(t, tf, tf) will fail
	assert.Equal(t, len(fb.formatters), 1)
	assert.Equal(t, reflect.ValueOf(fb.formatters[0]), reflect.ValueOf(f))
}

func TestRequired(t *testing.T) {
	fb := FieldBuilder{
		formatters: []FormatterFunc{},
		validators: []ValidatorFunc{},
	}
	assert.False(t, fb.required)
	fb.Required()
	assert.True(t, fb.required)
}

func TestEmpty(t *testing.T) {
	fb := FieldBuilder{
		formatters: []FormatterFunc{},
		validators: []ValidatorFunc{},
	}
	assert.Equal(t, fb.empty, nil)
	fb.Empty("bee")
	assert.Equal(t, fb.empty, "bee")
}

func TestBuild(t *testing.T) {
	f := FormatterFunc(func(s string) string { return s })

	fb := FieldBuilder{
		formatters: []FormatterFunc{f},
		validators: []ValidatorFunc{},
	}
	nf := fb.Build()
	assert.Equal(t, len(nf.formatters), 1)
	assert.Equal(t, reflect.ValueOf(nf.formatters[0]), reflect.ValueOf(f))
}
