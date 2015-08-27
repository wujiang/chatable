package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/wujiang/asapp"
)

type EnvelopesTestSuite struct {
	suite.Suite
}

func (s *EnvelopesTestSuite) SetupTest() {
	Init(testDB)
	createTables()
	newTestUsers()
	env, _ := asapp.NewEnvelope(1, 1, "test", 0)
	testStore.dbh.Insert(&env)
}

func (s *EnvelopesTestSuite) TearDownTest() {
	// dropTables()
	Exit()
}

func (s *EnvelopesTestSuite) TestGetByUserIDWithUserID() {

}

func TestEnvelopes(t *testing.T) {
	suite.Run(t, new(EnvelopesTestSuite))
}
