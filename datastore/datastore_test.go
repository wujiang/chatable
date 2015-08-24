package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	// TODO: add readme to setup test db
	testDB = "postgres://asapp@localhost:5432/asapp_test?sslmode=disable"
)

type DataStoreTestSuite struct {
	suite.Suite
}

// func (s *DataStoreTestSuite) SetupTest() {
// 	Init(testDB)
// 	CreateTables()
// }

// func (s *DataStoreTestSuite) TearDownTest() {
// 	DropTables()
// 	Exit()
// }

func (s *DataStoreTestSuite) TestInit() {
	Init(testDB)
	defer Exit()
	s.Nil(dbm.Db.Ping())
}

func (s *DataStoreTestSuite) TestExit() {
	Init(testDB)
	s.Nil(dbm.Db.Ping())
	Exit()
	s.NotNil(dbm.Db.Ping())
}

func TestDataStore(t *testing.T) {
	suite.Run(t, new(DataStoreTestSuite))
}
