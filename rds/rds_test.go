package rds

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/suite"
	"github.com/wujiang/chatable"
)

const (
	testRdsHost = "localhost:6379"
	testQueue   = "test:rds:queue"
	testQM      = "test:rds:qm"
)

var (
	testRdsPool = NewRdsPool(nil)
)

type RdsTestSuite struct {
	suite.Suite
}

func (rts *RdsTestSuite) SetupTest() {
	Init(testRdsHost)
}

func (rts *RdsTestSuite) TearDownTest() {
	conn := testRdsPool.pool.Get()
	defer conn.Close()
	conn.Do("DEL", testQueue, testQM)
	Exit()
}

func (rts *RdsTestSuite) TestEnqueue() {
	conn := testRdsPool.pool.Get()
	defer conn.Close()

	env := chatable.PublicEnvelope{
		Author:      "author",
		Recipient:   "recipient",
		Message:     "hello world",
		MessageType: chatable.MessageTypeText,
		CreatedAt:   time.Now().UTC(),
	}
	testRdsPool.Enqueue(testQueue, env)
	ct, err := redis.Int(conn.Do("LLEN", testQueue))
	rts.Nil(err)
	rts.Equal(1, ct)
}

func (rts *RdsTestSuite) TestDequeue() {
	env := chatable.PublicEnvelope{
		Author:      "author",
		Recipient:   "recipient",
		Message:     "hello world",
		MessageType: chatable.MessageTypeText,
		CreatedAt:   time.Now().UTC(),
	}
	testRdsPool.Enqueue(testQueue, env)
	e, err := testRdsPool.Dequeue(testQueue)
	rts.Nil(err)
	rts.Equal(env, e)
}

func (rts *RdsTestSuite) TestAddToQM() {
	conn := testRdsPool.pool.Get()
	defer conn.Close()

	rts.Nil(testRdsPool.AddToQM(testQM, "queue1"))
	ct, err := redis.Int(conn.Do("SCARD", testQM))
	rts.Equal(1, ct)
	rts.Nil(err)

	rts.Nil(testRdsPool.AddToQM(testQM, "queue1"))
	ct, err = redis.Int(conn.Do("SCARD", testQM))
	rts.Equal(1, ct)
	rts.Nil(err)

	rts.Nil(testRdsPool.AddToQM(testQM, "queue2"))
	ct, err = redis.Int(conn.Do("SCARD", testQM))
	rts.Equal(2, ct)
	rts.Nil(err)
}

func (rts *RdsTestSuite) TestQMMembers() {
	rts.Nil(testRdsPool.AddToQM(testQM, "queue1"))
	rts.Nil(testRdsPool.AddToQM(testQM, "queue1"))
	m, err := testRdsPool.QMMembers(testQM)
	rts.Equal([]string{"queue1"}, m)
	rts.Nil(err)

	rts.Nil(testRdsPool.AddToQM(testQM, "queue2"))
	m, err = testRdsPool.QMMembers(testQM)
	rts.Equal([]string{"queue2", "queue1"}, m)
	rts.Nil(err)
}

func (rts *RdsTestSuite) TestRemoveFromQM() {
	rts.Nil(testRdsPool.AddToQM(testQM, "queue1"))
	rts.Nil(testRdsPool.AddToQM(testQM, "queue2"))
	rts.Nil(testRdsPool.RemoveFromQM(testQM, "queue1"))
	m, err := testRdsPool.QMMembers(testQM)
	rts.Equal([]string{"queue2"}, m)
	rts.Nil(err)

	rts.Nil(testRdsPool.RemoveFromQM(testQM, "queue1"))
	m, err = testRdsPool.QMMembers(testQM)
	rts.Equal([]string{"queue2"}, m)
	rts.Nil(err)
}

func TestRds(t *testing.T) {
	suite.Run(t, new(RdsTestSuite))
}
