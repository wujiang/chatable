package asapp

import "time"

const (
	MessageTypeText = iota
	MessageTypePhoto
)

const (
	EnvelopesLimit = 1000
)

// Envelope represents a row in the envelopes table
type Envelope struct {
	ID          int      `db:"id"`
	UserID      int      `db:"user_id"`
	WithUserID  int      `db:"with_user_id"`
	IsIncoming  bool     `db:"is_incoming"`
	CreatedAt   NullTime `db:"created_at"`
	DeletedAt   NullTime `db:"deleted_at"`
	ReadAt      NullTime `db:"read_at"`
	Message     string   `db:"message"`
	MessageType int      `db:"message_type"`
}

func (env *Envelope) IsRead() bool {
	return env.ReadAt.Valid
}

func (env *Envelope) IsDeleted() bool {
	return env.DeletedAt.Valid
}

// NewEnvelope creates the incoming and outgoing envelopes. The first
// envelope is the envelope on the sender's side, and the second
// envelope is the envelope on the receipt's side.
func NewEnvelope(sender int, recipient int, msg string, t int) (*Envelope, *Envelope) {
	dt := NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	senderEnv := &Envelope{
		UserID:      sender,
		WithUserID:  recipient,
		IsIncoming:  false,
		CreatedAt:   dt,
		Message:     msg,
		MessageType: t,
	}
	recipientEnv := &Envelope{
		UserID:      recipient,
		WithUserID:  sender,
		IsIncoming:  true,
		CreatedAt:   dt,
		Message:     msg,
		MessageType: t,
	}
	return senderEnv, recipientEnv
}

// EnvelopeService defines the protocol for envelopes
type EnvelopeService interface {
	GetByUserIDWithUserID(uid int, withuid int, offset int) ([]*Envelope,
		error)
	Create(env *Envelope) error
	MarkDelete(env *Envelope) (int64, error)
	MarkRead(env *Envelope) (int64, error)
}