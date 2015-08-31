package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConnectionsTestSuite struct {
	suite.Suite
}

func (c *ConnectionsTestSuite) SetupTest() {
	Init(testDB)
	CreateTables()
	newTestUsers()
}

func (c *ConnectionsTestSuite) TearDownTest() {
	DropTables()
	Exit()
}

func (c *ConnectionsTestSuite) TestGetByUserID() {

}

func (c *ConnectionsTestSuite) TestDelete() {

}

func TestConnections(t *testing.T) {
	suite.Run(t, new(ConnectionsTestSuite))
}
