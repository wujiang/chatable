package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/wujiang/asapp"
)

type ThreadsTestSuite struct {
	suite.Suite
}

func (t *ThreadsTestSuite) SetupTest() {
	Init(testDB)
	createTables()
	newTestUsers()
}

func (t *ThreadsTestSuite) TearDownTest() {
	dropTables()
	Exit()
}

func (t *ThreadsTestSuite) TestUpsert() {
	thread := asapp.NewThread(1, 2, "recipient", "hello")
	ct, err := testStore.ThreadStore.Upsert(thread)
	t.Equal(int64(1), ct)
	t.Nil(err)

	var th asapp.Thread
	err = testStore.dbh.SelectOne(&th, "select * from threads where id = 1")
	t.Nil(err)
	t.Equal(thread.CreatedAt, th.CreatedAt)
	t.Equal(thread.LatestMessage, th.LatestMessage)

	newThread := asapp.NewThread(1, 2, "recipient", "what's going on")
	ct, err = testStore.ThreadStore.Upsert(newThread)
	t.Equal(int64(0), ct)
	t.Nil(err)

	var newTh asapp.Thread
	err = testStore.dbh.SelectOne(&newTh, "select * from threads where id = 1")
	t.Nil(err)
	t.Equal(newThread.CreatedAt, newTh.CreatedAt)
	t.Equal(newThread.LatestMessage, newTh.LatestMessage)
}

func (t *ThreadsTestSuite) TestGetByUserID() {
	thread := asapp.NewThread(1, 2, "recipient", "hello")
	err := testStore.dbh.Insert(thread)
	t.Nil(err)

	threads, err := testStore.ThreadStore.GetByUserID(1, 0)
	t.Equal(1, len(threads))
	t.Nil(err)
	t.Equal(thread.LatestMessage, threads[0].LatestMessage)
}

func TestThreads(t *testing.T) {
	suite.Run(t, new(ThreadsTestSuite))
}
