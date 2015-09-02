package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"gitlab.com/wujiang/asapp"
)

const (
	testDB = "postgres://asapp@localhost:5432/asapp_test?sslmode=disable"

	testUserPass        = "password"
	testUserFirst       = "first"
	testUserLast        = "last"
	testUsername        = "username"
	testUserEmail       = "first_last@example.com"
	testUserPhoneNumber = "1234567890"
	testUserIP          = "0.0.0.0"
)

func helperCreateUser() *asapp.User {
	u := asapp.NewUser(testUserFirst, testUserLast, testUsername,
		testUserPass, testUserEmail, testUserPhoneNumber,
		testUserIP)
	if err := store.UserStore.Create(u); err != nil {
		return nil
	} else {
		return u
	}
}

func helperAuthToken(host string) (asapp.PublicToken, error) {
	endpoint := strings.Join([]string{host, "/auth_token"}, "")
	payload := url.Values{
		"username": []string{testUsername},
		"password": []string{testUserPass},
	}
	resp, err := http.PostForm(endpoint, payload)
	if err != nil {
		return asapp.PublicToken{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data := struct {
		Data        []asapp.PublicToken `json:"data"`
		Total       int                 `json:"total"`
		CurrentPage int                 `json:"current_page"`
		PerPage     int                 `json:"per_page"`
	}{}
	if err = json.Unmarshal(body, &data); err != nil {
		return asapp.PublicToken{}, err
	}
	return data.Data[0], nil
}

func helperAuthHeader(host string, payload map[string]string) (string, error) {
	at, err := helperAuthToken(host)
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	for _, k := range payload {
		token.Claims[k] = payload[k]
	}
	token.Header["access_key"] = at.AccessKeyID
	t, err := token.SignedString([]byte(at.SecretAccessKey))
	return t, err
}
