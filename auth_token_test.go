package chatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNewAuthToken(t *testing.T) {
	tk1 := NewAuthToken(1, 1, []string{"all"})
	tk2 := NewAuthToken(1, 1, []string{"all"})
	assert.NotEqual(t, tk1.AccessKeyID, tk2.AccessKeyID)
	assert.NotEqual(t, tk1.SecretAccessKey, tk2.SecretAccessKey)
	assert.NotEqual(t, tk1.RefreshToken, tk2.RefreshToken)
}

type AuthTokenTestSuite struct {
	suite.Suite
	auth *AuthToken
}

func (a *AuthTokenTestSuite) SetupTest() {
	a.auth = NewAuthToken(1, 1, []string{"all"})
}

func (a *AuthTokenTestSuite) TestIsGood() {
	a.True(a.auth.IsGood())
	a.auth.IsActive = false
	a.False(a.auth.IsGood())
}

func (a *AuthTokenTestSuite) TestToPublicToken() {
	pub := a.auth.ToPublicToken()
	a.Equal(PublicToken{
		AccessKeyID:     a.auth.AccessKeyID,
		SecretAccessKey: a.auth.SecretAccessKey,
		RefreshToken:    a.auth.RefreshToken,
		CreatedAt:       a.auth.CreatedAt,
		ExpiresAt:       a.auth.ExpiresAt,
		ModifiedAt:      a.auth.ModifiedAt,
		IsRefreshable:   a.auth.IsRefreshable,
	}, *pub)
}

func TestAuthToken(t *testing.T) {
	suite.Run(t, new(AuthTokenTestSuite))
}
