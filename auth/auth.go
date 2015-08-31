package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"gitlab.com/wujiang/asapp"
	"gitlab.com/wujiang/asapp/datastore"
)

var (
	store = datastore.NewDataStore(nil)

	ErrUnauthenticated      = asapp.NewAuthenticationError("Unauthenticated")
	ErrAlreadyAuthenticated = asapp.NewAuthenticationError("Already authenticated")
	ErrUserNotFound         = asapp.NewAuthenticationError("User not found")
	ErrUnauthorized         = asapp.NewAuthenticationError("Unauthorized")
)

// protocol
// header: base64(access_token_id)
// claims: base64(payloads)
// signature: HMAC(header + claims)

// access_token generates a access_token_id and access_secret_key
// - access_token_id is like a username
// - access_secret_key is like a password
// Store them in postgres and send access_secret_key to clients when
// clients log in. Clients will use access_secret_key as the secret
// to sign JWT. Clients need to send the access_token_id with every request.
// Server use access_token_id to look for access_secret_key for a client.
//

type Token struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
}

const (
	TokenBearer = "Bearer"
	TokenJWT    = "JWT"
	ISS         = "asappd:JWT"
	ALG         = "HS256"
)

// keyfunc retrieves the secret access key from db using the access key id
// provided by user
var keyfunc = func(tk *jwt.Token) (interface{}, error) {
	accessKey := tk.Header["access_key"]
	if accessKey == nil {
		return nil, ErrUnauthenticated
	}
	at, err := store.AuthTokenStore.GetByAccessKeyID(accessKey.(string))
	if err != nil || !at.IsGood() {
		return nil, ErrUnauthenticated
	}
	// Add user to token's header so that we can add it to the request later
	user, err := store.UserStore.GetByID(at.UserID)
	if err != nil {
		return nil, ErrUnauthenticated
	}
	tk.Header["user"] = user
	tk.Header["auth"] = at
	return []byte(at.SecretAccessKey), nil
}

func TokenAuthenticate(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	token, err := jwt.ParseFromRequest(r, keyfunc)
	if err != nil || !token.Valid {
		return ErrUnauthenticated
	}
	context.Set(r, "user", token.Header["user"])
	context.Set(r, "auth", token.Header["auth"])
	return nil
}

func TokenUnAuthenticate(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	at := context.Get(r, "auth")
	if at == nil {
		return ErrUnauthenticated
	}
	authToken, ok := at.(*asapp.AuthToken)
	if !ok {
		return ErrUnauthenticated
	}
	authToken.IsActive = false
	if _, err := store.AuthTokenStore.Update(authToken); err != nil {
		return asapp.NewServerError(err.Error())
	}
	return nil
}
