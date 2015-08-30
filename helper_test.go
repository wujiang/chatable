package asapp

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHash(t *testing.T) {
	assert.True(t, "password" != GenerateHash("password"))
}

func TestCompareHash(t *testing.T) {
	assert.True(t, CompareHash(GenerateHash("password"), "password"))
}

func TesNewJSONResult(t *testing.T) {
	input := []struct{}{}
	result := NewJSONResult(input, 1)
	assert.Equal(t, *result, JSONResult{
		Data:        input,
		CurrentPage: 0,
		PerPage:     PerPage,
	})

	input2 := []struct {
		name  string
		email string
	}{
		struct {
			name  string
			email string
		}{
			name:  "test",
			email: "test@example.com",
		},
	}
	result = NewJSONResult(input2, 1)
	assert.Equal(t, *result, JSONResult{
		Data:        input2,
		CurrentPage: 1,
		PerPage:     PerPage,
	})
}

func TestGenerateRandomKey(t *testing.T) {
	for i := 0; i < 100; i++ {
		key := GenerateRandomKey()
		isMatched, err := regexp.MatchString("^[a-zA-Z0-9]+$", key)
		assert.Nil(t, err)
		assert.True(t, isMatched)
	}
}
