package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/wujiang/asapp"
	"gitlab.com/wujiang/asapp/datastore"
)

type envelopeTestData struct {
	Status      string                 `json:"status"`
	Error       asapp.JSONError        `json:"error"`
	Data        []asapp.PublicEnvelope `json:"data"`
	CurrentPage int                    `json:"current_page"`
	PerPage     int                    `json:"per_page"`
}

type EnvelopeTokenTestSuite struct {
	suite.Suite
	server *httptest.Server
	user   *asapp.User
}

func (s *EnvelopeTokenTestSuite) SetupTest() {
	s.server = httptest.NewServer(Handler())
	datastore.Init(testDB)
	datastore.CreateTables()
	s.user = helperCreateUser()
}

func (s *EnvelopeTokenTestSuite) TearDownTest() {
	s.server.Close()
	datastore.DropTables()
	datastore.Exit()
}

func (s *EnvelopeTokenTestSuite) TestGetEnvelopesUnauth() {
	endpoint := strings.Join([]string{s.server.URL, "thread", testUsername}, "/")

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	s.Nil(err)
	resp, err := client.Do(req)
	s.Nil(err)
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *EnvelopeTokenTestSuite) TestGetEnvelopes() {
	endpoint := strings.Join([]string{s.server.URL, "thread", testUsername}, "/")
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
	data := envelopeTestData{}
	err = json.Unmarshal(body, &data)
	s.Nil(err)
	s.Equal("success", data.Status)
}

func TestEnvelope(t *testing.T) {
	suite.Run(t, new(EnvelopeTokenTestSuite))
}
