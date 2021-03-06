package api

import (
	"net/http"
	"strings"

	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/forms"
)

var usersForm = forms.New()

func init() {
	usersForm.WithField("first_name", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader).
		WithValidators(forms.NameValidator()))
	usersForm.WithField("last_name", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader).
		WithValidators(forms.NameValidator()))
	usersForm.WithField("username", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader).
		WithValidators(forms.UsernameValidator()))
	usersForm.WithField("password", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader).
		WithValidators(forms.PasswordValidator()))
	usersForm.WithField("email", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader).
		WithValidators(forms.EmailValidator()))
	usersForm.WithField("phone", new(forms.FieldBuilder).
		Required().
		Loader(forms.StringLoader).
		WithValidators(forms.PhoneNumberValidator()))
}

func serveCreateUser(w http.ResponseWriter, r *http.Request) chatable.CompoundError {
	if err := r.ParseForm(); err != nil {
		return chatable.NewServerError(err.Error())
	}
	valid := usersForm.Valid(r.PostForm)
	if !valid {
		return usersForm.ConsolidateErrors()
	}
	u := chatable.NewUser(usersForm.Values["first_name"].(string),
		usersForm.Values["last_name"].(string),
		usersForm.Values["username"].(string),
		usersForm.Values["password"].(string),
		usersForm.Values["email"].(string),
		usersForm.Values["phone"].(string),
		r.RemoteAddr)
	err := store.UserStore.Create(u)
	// TODO: refine the error message
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "violates") {
			return chatable.NewUserError("Some fileds are not unique")
		}
		return chatable.NewServerError(msg)
	}
	// create a auth token
	token, err := createNewAuthToken(w, r, u)
	if err != nil {
		return chatable.NewServerError(err.Error())
	}
	data := []chatable.UserWithToken{
		chatable.UserWithToken{
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Username:    u.Username,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Token:       *token,
		},
	}
	return writeJSON(w, chatable.NewJSONResult(data, 1))
}
