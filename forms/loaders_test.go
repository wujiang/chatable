package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringLoader(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]string{
		"":         "",
		"  ":       "",
		"  bee ":   "bee",
		"bee app ": "bee app",
	}
	for k, v := range tests {
		cv, err := StringLoader(k)
		assert.Equal(v, cv)
		assert.Nil(err)
	}
}

func TestIntLoader(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]int{
		"1":   1,
		"123": 123,
	}
	for k, v := range tests {
		cv, err := IntLoader(k)
		assert.Equal(v, cv)
		assert.Equal(err, nil)
	}
	failedTests := []string{"", "1s"}
	for _, v := range failedTests {
		cv, err := IntLoader(v)
		assert.Nil(cv)
		assert.Equal(err, ErrNotANumber)
	}
}
