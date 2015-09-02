package rds

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/suite"
	"gitlab.com/wujiang/asapp"
)

const (
	testRdsHost = "localhost:6379"
	testQueue   = "test:rds:queue"
	testQM      = "test:rds:qm"
)

var (
	testRdsConn = NewRdsConn(nil)
)

type RdsTestSuite struct {
	suite.Suite
}

func (rts *RdsTestSuite) SetupTest() {
	Init(testRdsHost)
}

func (rts *RdsTestSuite) TearDownTest() {
	(*testRdsConn.conn).Do("DEL", testQueue, testQM)
	Exit()
}

func (rts *RdsTestSuite) TestEnqueue() {
	env := asapp.PublicEnvelope{
		Author:      "author",
		Recipient:   "recipient",
		Message:     "hello world",
		MessageType: asapp.MessageTypeText,
		CreatedAt:   time.Now().UTC(),
	}
	testRdsConn.Enqueue(testQueue, env)
	ct, err := redis.Int((*testRdsConn.conn).Do("LLEN", testQueue))
	rts.Nil(err)
	rts.Equal(1, ct)
}

func (rts *RdsTestSuite) TestDequeue() {
	env := asapp.PublicEnvelope{
		Author:      "author",
		Recipient:   "recipient",
		Message:     "hello world",
		MessageType: asapp.MessageTypeText,
		CreatedAt:   time.Now().UTC(),
	}
	testRdsConn.Enqueue(testQueue, env)
	e, err := testRdsConn.Dequeue(testQueue)
	rts.Nil(err)
	rts.Equal(env, e)
}

func (rts *RdsTestSuite) TestAddToQM() {
	rts.Nil(testRdsConn.AddToQM(testQM, "queue1"))
	ct, err := redis.Int((*testRdsConn.conn).Do("SCARD", testQM))
	rts.Equal(1, ct)
	rts.Nil(err)

	rts.Nil(testRdsConn.AddToQM(testQM, "queue1"))
	ct, err = redis.Int((*testRdsConn.conn).Do("SCARD", testQM))
	rts.Equal(1, ct)
	rts.Nil(err)

	rts.Nil(testRdsConn.AddToQM(testQM, "queue2"))
	ct, err = redis.Int((*testRdsConn.conn).Do("SCARD", testQM))
	rts.Equal(2, ct)
	rts.Nil(err)
}

func (rts *RdsTestSuite) TestQMMembers() {
	rts.Nil(testRdsConn.AddToQM(testQM, "queue1"))
	rts.Nil(testRdsConn.AddToQM(testQM, "queue1"))
	m, err := testRdsConn.QMMembers(testQM)
	rts.Equal([]string{"queue1"}, m)
	rts.Nil(err)

	rts.Nil(testRdsConn.AddToQM(testQM, "queue2"))
	m, err = testRdsConn.QMMembers(testQM)
	rts.Equal([]string{"queue1", "queue2"}, m)
	rts.Nil(err)
}

func (rts *RdsTestSuite) TestRemoveFromQM() {
	rts.Nil(testRdsConn.AddToQM(testQM, "queue1"))
	rts.Nil(testRdsConn.AddToQM(testQM, "queue2"))
	rts.Nil(testRdsConn.RemoveFromQM(testQM, "queue1"))
	m, err := testRdsConn.QMMembers(testQM)
	rts.Equal([]string{"queue2"}, m)
	rts.Nil(err)

	rts.Nil(testRdsConn.RemoveFromQM(testQM, "queue1"))
	m, err = testRdsConn.QMMembers(testQM)
	rts.Equal([]string{"queue2"}, m)
	rts.Nil(err)
}

func TestRds(t *testing.T) {
	suite.Run(t, new(RdsTestSuite))
}
