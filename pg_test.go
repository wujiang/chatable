package asapp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullTimeScan(t *testing.T) {
	var nt NullTime
	err := nt.Scan(nil)
	assert.Nil(t, err)
	assert.False(t, nt.Valid)

	err = nt.Scan([]byte("2006-01-02T15:04:05.999999999Z"))
	assert.Nil(t, err)
	assert.True(t, nt.Valid)
	assert.Equal(t, time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.UTC),
		nt.Time)

	err = nt.Scan("2006-01-02T15:04:05.999999999Z")
	assert.NotNil(t, err)
}

func TestNullTimeValue(t *testing.T) {
	now := time.Now().UTC()
	nt := NullTime{
		Time:  now,
		Valid: true,
	}
	v, err := nt.Value()
	assert.Nil(t, err)
	assert.Equal(t, now, v.(time.Time))

	nt = NullTime{
		Valid: false,
	}
	v, err = nt.Value()
	assert.Nil(t, err)
	assert.Nil(t, v)
}
