package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/datastore"
)

type UsersTestSuite struct {
	suite.Suite
	server *httptest.Server
	user   *chatable.User
}

type usersTestData struct {
	Status      string                   `json:"status"`
	Error       chatable.JSONError       `json:"error"`
	Data        []chatable.UserWithToken `json:"data"`
	CurrentPage int                      `json:"current_page"`
	PerPage     int                      `json:"per_page"`
}

func (s *UsersTestSuite) SetupTest() {
	s.server = httptest.NewServer(Handler())
	datastore.Init(testDB)
	datastore.CreateTables()
	s.user = helperCreateUser()
}

func (s *UsersTestSuite) TearDownTest() {
	s.server.Close()
	datastore.DropTables()
	datastore.Exit()
}

func (s *UsersTestSuite) TestCreateUsersMissing() {
	endpoint := strings.Join([]string{s.server.URL, "/register"}, "")
	payload := url.Values{}
	payload.Set("first_name", "Test")
	payload.Set("last_name", "Last")
	resp, err := http.PostForm(endpoint, payload)
	s.Nil(err)
	s.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var data = usersTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)

	s.Equal(usersTestData{
		Status: "fail",
		Error: chatable.JSONError{
			Code:    http.StatusBadRequest,
			Message: "Form errors",
			Errors: chatable.ErrorDetails{
				"username": "is required",
				"email":    "is required",
				"password": "is required",
				"phone":    "is required",
			},
		},
		CurrentPage: 0,
		PerPage:     chatable.PerPage,
		Data:        []chatable.UserWithToken{},
	}, data)

}

func (s *UsersTestSuite) TestCreateUsersFormat() {
	endpoint := strings.Join([]string{s.server.URL, "/register"}, "")
	payload := url.Values{}
	payload.Set("username", "username")
	payload.Set("first_name", "Test")
	payload.Set("last_name", "Last")
	payload.Set("email", "Test_Last")
	payload.Set("phone", "135792468")
	payload.Set("password", "pass")

	resp, err := http.PostForm(endpoint, payload)
	s.Nil(err)
	s.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var data = usersTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)

	s.Equal(usersTestData{
		Status: "fail",
		Error: chatable.JSONError{
			Code:    http.StatusBadRequest,
			Message: "Form errors",
			Errors: chatable.ErrorDetails{
				"email":    "please enter a valid email address",
				"password": "please enter a password with at least 10 characters",
			},
		},
		CurrentPage: 0,
		PerPage:     chatable.PerPage,
		Data:        []chatable.UserWithToken{},
	}, data)
}

func (s *UsersTestSuite) TestCreateUsers() {
	endpoint := strings.Join([]string{s.server.URL, "/register"}, "")
	payload := url.Values{}
	payload.Set("username", "username")
	payload.Set("first_name", "Test")
	payload.Set("last_name", "Last")
	payload.Set("email", "Test_Last@example.com")
	payload.Set("phone", "135792468")
	payload.Set("password", "password123")

	resp, err := http.PostForm(endpoint, payload)
	s.Nil(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var data = usersTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)

	s.Equal("success", data.Status)
	s.Equal(chatable.JSONError{
		Code:    http.StatusOK,
		Message: "",
		Errors:  nil,
	}, data.Error)
	s.Equal(1, data.CurrentPage)
	s.Equal("Test", data.Data[0].FirstName)
	s.Equal("Last", data.Data[0].LastName)
	s.Equal("135792468", data.Data[0].PhoneNumber)
	s.Equal("Test_Last@example.com", data.Data[0].Email)
	s.NotNil(data.Data[0].Token)
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
