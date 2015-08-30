package api

import (
	"net/http"
	"strconv"

	"gitlab.com/wujiang/asapp"
	"gitlab.com/wujiang/asapp/auth"
	"gitlab.com/wujiang/asapp/forms"
)

var (
	authTokenForm = forms.New()
)

func init() {
	authTokenForm.WithField("email", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader))
	authTokenForm.WithField("password", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader))
}

// createNewAuthToken creates a new auth token
func createNewAuthToken(w http.ResponseWriter, r *http.Request, u *asapp.User) (*asapp.PublicToken, asapp.CompoundError) {
	// create a new token for the user
	// client_id is on the header
	clientID := r.Header.Get("ClientID")
	cid, err := strconv.Atoi(clientID)
	if err != nil {
		cid = -1
	}
	at := asapp.NewAuthToken(u.ID, cid, asapp.StringSlice{"all"})
	if err = store.AuthTokenStore.Create(at); err != nil {
		return nil, asapp.NewServerError(err.Error())
	}
	return at.ToPublicToken(), nil
}

func serveCreateAuthToken(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	if err := r.ParseForm(); err != nil {
		return asapp.NewServerError(err.Error())
	}
	valid := authTokenForm.Valid(r.PostForm)
	if !valid {
		return authTokenForm.ConsolidateErrors()
	}
	u, err := store.UserStore.GetByEmail(authTokenForm.Values["email"].(string))
	if err != nil {
		return auth.ErrUnauthenticated
	}
	if !asapp.CompareHash(u.Password, authTokenForm.Values["password"].(string)) {
		return auth.ErrUnauthenticated
	}

	token, cerr := createNewAuthToken(w, r, u)
	if cerr != nil {
		return cerr
	}
	return writeJSON(w, asapp.NewJSONResult([]*asapp.PublicToken{token}, 1))
}

func serveDeactivateAuthToken(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	return auth.TokenUnAuthenticate(w, r)
}
