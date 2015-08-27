package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DataStoreTestSuite struct {
	suite.Suite
}

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
