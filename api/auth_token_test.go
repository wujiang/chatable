package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/datastore"

	"github.com/stretchr/testify/suite"
)

type authTestData struct {
	Status      string                 `json:"status"`
	Error       chatable.JSONError     `json:"error"`
	Data        []chatable.PublicToken `json:"data"`
	CurrentPage int                    `json:"current_page"`
	PerPage     int                    `json:"per_page"`
}

type AuthTokenTestSuite struct {
	suite.Suite
	server *httptest.Server
	user   *chatable.User
}

func (s *AuthTokenTestSuite) SetupTest() {
	s.server = httptest.NewServer(Handler())
	datastore.Init(testDB)
	datastore.CreateTables()
	s.user = helperCreateUser()
}

func (s *AuthTokenTestSuite) TearDownTest() {
	s.server.Close()
	datastore.DropTables()
	datastore.Exit()
}

func (s *AuthTokenTestSuite) TestCreateAuthToken() {
	endpoint := strings.Join([]string{s.server.URL, "auth_token"}, "/")
	payload := url.Values{}
	payload.Set("username", testUsername)
	payload.Set("password", testUserPass)
	resp, err := http.PostForm(endpoint, payload)
	s.Nil(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	data := authTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)

	s.Equal("success", data.Status)
	s.Equal(chatable.JSONError{
		Code:    http.StatusOK,
		Message: "",
		Errors:  chatable.ErrorDetails(nil),
	}, data.Error)
	s.Equal(1, len(data.Data))
}

func (s *AuthTokenTestSuite) TestCreateAuthTokenInvalid() {
	endpoint := strings.Join([]string{s.server.URL, "auth_token"}, "/")
	payload := url.Values{}
	payload.Set("username", testUsername)
	payload.Set("password", "wrong")
	resp, err := http.PostForm(endpoint, payload)
	s.Nil(err)
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	data := authTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)

	s.Equal("fail", data.Status)
	s.Equal(chatable.JSONError{
		Code:    http.StatusUnauthorized,
		Message: "Authentication errors",
		Errors: chatable.ErrorDetails{
			"error": "Unauthenticated",
		},
	}, data.Error)
	s.Equal(0, len(data.Data))
}

func (s *AuthTokenTestSuite) TestDeactivateAuthToken() {
	endpoint := strings.Join([]string{s.server.URL, "auth_token"}, "/")
	header, err := helperAuthHeader(s.server.URL, map[string]string{})
	s.Nil(err)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", endpoint, nil)
	s.Nil(err)
	header = strings.Join([]string{"BEARER", header}, " ")
	req.Header.Add("Authorization", header)
	resp, err := client.Do(req)
	s.Nil(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data := authTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)
	s.Equal("success", data.Status)
}

func TestAuthToken(t *testing.T) {
	suite.Run(t, new(AuthTokenTestSuite))
}
