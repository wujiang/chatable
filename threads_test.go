package chatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ThreadTestSuite struct {
	suite.Suite
	outgoing *Thread
	incoming *Thread
}

func (e *ThreadTestSuite) SetupTest() {
	th1, th2 := NewThread(1, 2, "sender", "hello")
	e.outgoing = th1
	e.incoming = th2
}

func (e *ThreadTestSuite) TestToPublic() {
	t := e.outgoing.ToPublic()
	e.Equal(PublicThread{
		AuthorUsername: e.outgoing.AuthorUsername,
		CreatedAt:      e.outgoing.CreatedAt.Time,
		LatestMessage:  e.outgoing.LatestMessage,
	}, *t)
}

func TestThread(t *testing.T) {
	suite.Run(t, new(ThreadTestSuite))
}

func TestNewThread(t *testing.T) {
	out, in := NewThread(1, 2, "sender", "hello")
	assert.Equal(t, out.CreatedAt, in.CreatedAt)
	assert.Equal(t, out.LatestMessage, in.LatestMessage)
	assert.Equal(t, out.AuthorUsername, in.AuthorUsername)
}
