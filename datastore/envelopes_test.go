package datastore

import "github.com/stretchr/testify/suite"

type EnvelopesTestSuite struct {
	suite.Suite
}

func (s *EnvelopesTestSuite) SetupTest() {
	Init(testDB)
	CreateTables()
}

func (s *EnvelopesTestSuite) TearDownTest() {
	DropTables()
	Exit()
}
