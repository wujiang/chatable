package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UsersTestSuite struct {
	suite.Suite
}

func (s *UsersTestSuite) SetupTest() {
	Init(testDB)
	createTables()
	newTestUsers()
}

func (s *UsersTestSuite) TearDownTest() {
	dropTables()
	Exit()
}

func (s *UsersTestSuite) TestGetByUsername() {
	user, err := testStore.UserStore.GetByUsername(testSenderUname)
	s.Equal(testSender.Username, user.Username)
	s.Equal(testSender.Email, user.Email)
	s.Equal(testSender.PhoneNumber, user.PhoneNumber)
	s.Equal(testSender.Password, user.Password)
	s.Nil(err)
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
