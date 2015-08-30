package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapFormatter(t *testing.T) {
	f := CapFormatter
	assert.Equal(t, f(""), "")
	assert.Equal(t, f("BEE"), "Bee")
	assert.Equal(t, f("bee"), "Bee")
	assert.Equal(t, f("bEe"), "Bee")
}
