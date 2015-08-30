package forms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailValidator(t *testing.T) {
	vf := EmailValidator()
	assert.Equal(t, vf(""), ErrEmailFormat)
	assert.Equal(t, vf("name"), ErrEmailFormat)
	assert.Equal(t, vf("name@beeapp"), ErrEmailFormat)
	assert.Equal(t, vf("name@.beeapp.com"), ErrEmailFormat)
	assert.Equal(t, vf("name@beeapp.io"), nil)
	assert.Equal(t, vf("#%name@beeapp.io"), nil)
	assert.Equal(t, vf("#%name@test.beeapp.io"), nil)
}

func TestPasswordValidator(t *testing.T) {
	vf := PasswordValidator()
	assert.Equal(t, vf("  "), ErrPasswordLen)
	assert.Equal(t, vf("password"), ErrPasswordLen)
	assert.Equal(t, vf("password password"), ErrPasswordSpace)
	assert.Equal(t, vf("passwordpassword"), nil)
}
