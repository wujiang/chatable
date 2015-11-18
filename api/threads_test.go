package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/datastore"
)

type threadTestData struct {
	Status      string                  `json:"status"`
	Error       chatable.JSONError      `json:"error"`
	Data        []chatable.PublicThread `json:"data"`
	CurrentPage int                     `json:"current_page"`
	PerPage     int                     `json:"per_page"`
}

type ThreadTokenTestSuite struct {
	suite.Suite
	server *httptest.Server
	user   *chatable.User
}

func (s *ThreadTokenTestSuite) SetupTest() {
	s.server = httptest.NewServer(Handler())
	datastore.Init(testDB)
	datastore.CreateTables()
	s.user = helperCreateUser()
}

func (s *ThreadTokenTestSuite) TearDownTest() {
	s.server.Close()
	datastore.DropTables()
	datastore.Exit()
}

func (s *ThreadTokenTestSuite) TestGetThreadsUnauth() {
	endpoint := strings.Join([]string{s.server.URL, "inbox"}, "/")

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	s.Nil(err)
	resp, err := client.Do(req)
	s.Nil(err)

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *ThreadTokenTestSuite) TestGetThreads() {
	endpoint := strings.Join([]string{s.server.URL, "inbox"}, "/")
	header, err := helperAuthHeader(s.server.URL, map[string]string{})
	s.Nil(err)

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	s.Nil(err)
	header = strings.Join([]string{"BEARER", header}, " ")
	req.Header.Add("Authorization", header)
	resp, err := client.Do(req)
	s.Nil(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data := threadTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)
	s.Equal("success", data.Status)
}

func TestThread(t *testing.T) {
	suite.Run(t, new(ThreadTokenTestSuite))
}
