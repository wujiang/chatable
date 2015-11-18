package chatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EnvelopeTestSuite struct {
	suite.Suite
	outgoing *Envelope
	incoming *Envelope
}

func (e *EnvelopeTestSuite) SetupTest() {
	env1, env2 := NewEnvelope(1, 2, "hello", MessageTypeText)
	e.outgoing = env1
	e.incoming = env2
}

func (e *EnvelopeTestSuite) TestIsRead() {
	e.False(e.outgoing.IsRead())
	e.False(e.outgoing.IsRead())
}

func (e *EnvelopeTestSuite) TestIsDeleted() {
	e.False(e.outgoing.IsDeleted())
	e.False(e.outgoing.IsDeleted())
}

func TestEnvelope(t *testing.T) {
	suite.Run(t, new(EnvelopeTestSuite))
}

func TestNewEnvelope(t *testing.T) {
	out, in := NewEnvelope(1, 2, "hello", MessageTypeText)
	assert.Equal(t, out.CreatedAt, in.CreatedAt)
	assert.Equal(t, out.Message, in.Message)
	assert.Equal(t, out.MessageType, in.MessageType)
	assert.Equal(t, out.UserID, 1)
	assert.Equal(t, out.WithUserID, 2)
	assert.Equal(t, in.WithUserID, 1)
	assert.Equal(t, in.UserID, 2)
}
