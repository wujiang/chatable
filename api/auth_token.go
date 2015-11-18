package api

import (
	"net/http"
	"strconv"

	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/auth"
	"github.com/wujiang/chatable/forms"
)

var (
	authTokenForm = forms.New()
)

func init() {
	authTokenForm.WithField("username", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader))
	authTokenForm.WithField("password", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader))
}

// createNewAuthToken creates a new auth token
func createNewAuthToken(w http.ResponseWriter, r *http.Request, u *chatable.User) (*chatable.PublicToken, chatable.CompoundError) {
	// create a new token for the user
	// client_id is on the header
	clientID := r.Header.Get("ClientID")
	cid, err := strconv.Atoi(clientID)
	if err != nil {
		cid = -1
	}
	at := chatable.NewAuthToken(u.ID, cid, chatable.StringSlice{"all"})
	if err = store.AuthTokenStore.Create(at); err != nil {
		return nil, chatable.NewServerError(err.Error())
	}
	return at.ToPublicToken(), nil
}

func serveCreateAuthToken(w http.ResponseWriter, r *http.Request) chatable.CompoundError {
	if err := r.ParseForm(); err != nil {
		return chatable.NewServerError(err.Error())
	}
	valid := authTokenForm.Valid(r.PostForm)
	if !valid {
		return authTokenForm.ConsolidateErrors()
	}
	u, err := store.UserStore.GetByUsername(authTokenForm.Values["username"].(string))
	if err != nil {
		return auth.ErrUnauthenticated
	}
	if !chatable.CompareHash(u.Password, authTokenForm.Values["password"].(string)) {
		return auth.ErrUnauthenticated
	}

	token, cerr := createNewAuthToken(w, r, u)
	if cerr != nil {
		return cerr
	}
	return writeJSON(w, chatable.NewJSONResult([]*chatable.PublicToken{token}, 1))
}

func serveDeactivateAuthToken(w http.ResponseWriter, r *http.Request) chatable.CompoundError {
	if err := auth.TokenUnAuthenticate(w, r); err != nil {
		return err
	}
	return writeJSON(w, chatable.NewJSONResult([]*chatable.PublicToken{}, 1))
}
