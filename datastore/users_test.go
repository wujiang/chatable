package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/wujiang/asapp"
)

type UsersTestSuite struct {
	suite.Suite
}

func (s *UsersTestSuite) SetupTest() {
	Init(testDB)
	CreateTables()
	newTestUsers()
}

func (s *UsersTestSuite) TearDownTest() {
	DropTables()
	Exit()
}

func (s *UsersTestSuite) TestGetByID() {
	user, err := testStore.UserStore.GetByID(1)
	s.Equal(testSender.Username, user.Username)
	s.Equal(testSender.Email, user.Email)
	s.Equal(testSender.PhoneNumber, user.PhoneNumber)
	s.Equal(testSender.Password, user.Password)
	s.Nil(err)
}

func (s *UsersTestSuite) TestGetByIDs() {
	users, err := testStore.UserStore.GetByIDs(1, 2)
	s.Equal(2, len(users))
	s.Nil(err)
	// s.Equal(testSender.Username, user.Username)
	// s.Equal(testSender.Email, user.Email)
	// s.Equal(testSender.PhoneNumber, user.PhoneNumber)
	// s.Equal(testSender.Password, user.Password)
	s.Nil(err)
}

func (s *UsersTestSuite) TestGetByUsername() {
	user, err := testStore.UserStore.GetByUsername(testSenderUname)
	s.Equal(testSender.Username, user.Username)
	s.Equal(testSender.Email, user.Email)
	s.Equal(testSender.PhoneNumber, user.PhoneNumber)
	s.Equal(testSender.Password, user.Password)
	s.Nil(err)
}

func (s *UsersTestSuite) TestCreate() {
	user := asapp.NewUser("test", "last", "username", "password123",
		"test@last.com", "1357902468", "0.0.0.0")
	s.Nil(testStore.UserStore.Create(user))
	u, err := testStore.UserStore.GetByUsername("username")
	s.Nil(err)
	s.Equal("test", u.FirstName)
	s.Equal("username", u.Username)

	// duplicates
	s.NotNil(testStore.UserStore.Create(user))
	ct, err := testStore.dbh.Delete(user)
	s.Equal(int64(1), ct)
	s.Nil(err)
}

func (s *UsersTestSuite) TestUpdate() {
	u := testSender
	u.Email = "changed@send.com"
	ct, err := testStore.UserStore.Update(u)
	s.Equal(int64(1), ct)
	s.Nil(err)

	user, err := testStore.UserStore.GetByUsername(testSenderUname)
	s.Equal(testSender.Username, user.Username)
	s.Equal(testSender.Email, "changed@send.com")
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
