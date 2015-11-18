package datastore

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ThreadsTestSuite struct {
	suite.Suite
}

func (t *ThreadsTestSuite) SetupTest() {
	Init(testDB)
	CreateTables()
	newTestUsers()
}

func (t *ThreadsTestSuite) TearDownTest() {
	DropTables()
	Exit()
}

func (t *ThreadsTestSuite) TestUpsert() {
	thread1, thread2 := chatable.NewThread(1, 2, "recipient", "hello")
	ct, err := testStore.ThreadStore.Upsert(thread1)
	t.Equal(int64(1), ct)
	t.Nil(err)
	ct, err = testStore.ThreadStore.Upsert(thread2)
	t.Nil(err)

	var th chatable.Thread
	err = testStore.dbh.SelectOne(&th, "select * from threads where id = 1")
	t.Nil(err)
	t.Equal(thread1.CreatedAt, th.CreatedAt)
	t.Equal(thread1.LatestMessage, th.LatestMessage)

	newThread1, _ := chatable.NewThread(1, 2, "recipient", "what's going on")
	ct, err = testStore.ThreadStore.Upsert(newThread1)
	t.Equal(int64(0), ct)
	t.Nil(err)

	var newTh chatable.Thread
	err = testStore.dbh.SelectOne(&newTh, "select * from threads where id = 1")
	t.Nil(err)
	t.Equal(newThread1.CreatedAt, newTh.CreatedAt)
	t.Equal(newThread1.LatestMessage, newTh.LatestMessage)
}

func (t *ThreadsTestSuite) TestGetByUserID() {
	thread1, thread2 := chatable.NewThread(1, 2, "recipient", "hello")
	err := testStore.dbh.Insert(thread1)
	t.Nil(err)
	err = testStore.dbh.Insert(thread2)
	t.Nil(err)

	threads, err := testStore.ThreadStore.GetByUserID(1, 0)
	t.Equal(1, len(threads))
	t.Nil(err)
	t.Equal(thread1.LatestMessage, threads[0].LatestMessage)
}

func TestThreads(t *testing.T) {
	suite.Run(t, new(ThreadsTestSuite))
}
