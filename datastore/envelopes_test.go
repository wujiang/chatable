package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvelopesTestSuite struct {
	suite.Suite
}

func (s *EnvelopesTestSuite) SetupTest() {
	Init(testDB)
	CreateTables()
	newTestUsers()
}

func (s *EnvelopesTestSuite) TearDownTest() {
	DropTables()
	Exit()
}

func (s *EnvelopesTestSuite) TestCreate() {
	sEnv, rEnv := chatable.NewEnvelope(1, 2, "hello", chatable.MessageTypeText)
	s.Nil(testStore.EnvelopeStore.Create(sEnv))
	s.Nil(testStore.EnvelopeStore.Create(rEnv))
}

func (s *EnvelopesTestSuite) TestGetByUserIDWithUserID() {
	sEnv, rEnv := chatable.NewEnvelope(1, 2, "hello", chatable.MessageTypeText)
	s.Nil(testStore.EnvelopeStore.Create(sEnv))
	s.Nil(testStore.EnvelopeStore.Create(rEnv))
	envs, err := testStore.EnvelopeStore.GetByUserIDWithUserID(1, 2, 0)
	s.Nil(err)
	s.Equal(1, len(envs))
	env := envs[0]
	s.Equal(sEnv.CreatedAt, env.CreatedAt)
	s.Equal(sEnv.UserID, env.UserID)
	s.Equal(sEnv.WithUserID, env.WithUserID)
	s.Equal(sEnv.DeletedAt, env.DeletedAt)
	s.Equal(sEnv.ReadAt, env.ReadAt)
	s.Equal(sEnv.Message, env.Message)
	s.Equal(sEnv.MessageType, env.MessageType)
}

func (s *EnvelopesTestSuite) TestMarkDelete() {
	sEnv, _ := chatable.NewEnvelope(1, 2, "hello", chatable.MessageTypeText)
	s.Nil(testStore.EnvelopeStore.Create(sEnv))

	env := chatable.Envelope{}
	err := testStore.dbh.SelectOne(&env, "select * from envelopes where id = 1")
	s.False(env.DeletedAt.Valid)

	ct, err := testStore.EnvelopeStore.MarkDelete(sEnv)
	s.Equal(int64(1), ct)
	s.Nil(err)

	err = testStore.dbh.SelectOne(&env, "select * from envelopes where id = 1")
	s.True(env.DeletedAt.Valid)
}

func (s *EnvelopesTestSuite) TestMarkRead() {
	sEnv, _ := chatable.NewEnvelope(1, 2, "hello", chatable.MessageTypeText)
	s.Nil(testStore.EnvelopeStore.Create(sEnv))

	env := chatable.Envelope{}
	err := testStore.dbh.SelectOne(&env, "select * from envelopes where id = 1")
	s.False(env.ReadAt.Valid)

	ct, err := testStore.EnvelopeStore.MarkRead(sEnv)
	s.Equal(int64(1), ct)
	s.Nil(err)

	err = testStore.dbh.SelectOne(&env, "select * from envelopes where id = 1")
	s.True(env.ReadAt.Valid)
}

func TestEnvelopes(t *testing.T) {
	suite.Run(t, new(EnvelopesTestSuite))
}
