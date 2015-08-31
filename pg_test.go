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

	err = nt.Scan(time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.UTC))
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

func TestStringSliceScan(t *testing.T) {
	var slice StringSlice

	err := slice.Scan([]byte(`{"12",45,"abc,\\\"d\\ef\\\\"}`))

	if err != nil {
		t.Errorf("Could not scan array, %v", err)
		return
	}

	if slice[0] != "12" || slice[2] != `abc,\"d\ef\\` {
		t.Errorf("Did not get expected slice contents")
	}
}

func TestStringSliceDbValue(t *testing.T) {
	slice := StringSlice([]string{`as"f\df`, "43", "}adsf"})

	val, err := slice.Value()
	if err != nil {
		t.Errorf("Could not convert to db string")
	}

	if str, ok := val.(string); ok {
		if `{"as\"f\\\df","43","}adsf"}` != str {
			t.Errorf("db value expecting %s got %s", `{"as\"f\\\df","43","}adsf"}`, str)
		}
	} else {
		t.Errorf("Could not convert %v to string for comparison", val)
	}
}
