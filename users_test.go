package chatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	u := NewUser("first", "last", "uname", "password", "hello",
		"12345", "0.0.0.0")
	assert.NotEqual(t, u.Password, "password")
}
