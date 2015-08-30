package asapp

import "time"

const (
	maxDuration = time.Hour * time.Duration(24*30)
)

type AuthToken struct {
	ID              int         `db:"id"`
	AccessKeyID     string      `db:"access_key_id"`
	SecretAccessKey string      `db:"secret_access_key"`
	RefreshToken    string      `db:"refresh_token"`
	CreatedAt       time.Time   `db:"created_at"`
	ExpiresAt       time.Time   `db:"expires_at"`
	ModifiedAt      time.Time   `db:"modified_at"`
	UserID          int         `db:"user_id"`
	ClientID        int         `db:"client_id"`
	IsActive        bool        `db:"is_active"`
	IsRefreshable   bool        `db:"is_refreshable"`
	Scope           StringSlice `db:"scope"`
}

func (at *AuthToken) IsGood() bool {
	return at.IsActive && at.ExpiresAt.After(time.Now().UTC())
}

func (at *AuthToken) ToPublicToken() *PublicToken {
	return &PublicToken{
		AccessKeyID:     at.AccessKeyID,
		SecretAccessKey: at.SecretAccessKey,
		RefreshToken:    at.RefreshToken,
		CreatedAt:       at.CreatedAt,
		ExpiresAt:       at.ExpiresAt,
		ModifiedAt:      at.ModifiedAt,
		IsRefreshable:   at.IsRefreshable,
	}
}

type AuthTokenService interface {
	GetByAccessKeyID(key string) (*AuthToken, error)
	Create(at *AuthToken) error
	Update(at *AuthToken) (int64, error)
}

func NewAuthToken(uid int, cid int, scope StringSlice) *AuthToken {
	return &AuthToken{
		AccessKeyID:     GenerateRandomKey(),
		SecretAccessKey: GenerateRandomKey(),
		RefreshToken:    GenerateRandomKey(),
		CreatedAt:       time.Now().UTC(),
		ExpiresAt:       time.Now().UTC().Add(maxDuration),
		ModifiedAt:      time.Now().UTC(),
		UserID:          uid,
		ClientID:        cid,
		IsActive:        true,
		IsRefreshable:   true,
		Scope:           scope,
	}
}

type PublicToken struct {
	AccessKeyID     string    `json:"access_key_id"`
	SecretAccessKey string    `json:"secret_access_key"`
	RefreshToken    string    `json:"refresh_token"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	ModifiedAt      time.Time `json:"modified_at"`
	IsRefreshable   bool      `json:"is_refreshable"`
}
