package asapp

import "time"

// Thread represents a row in the threads table
type Thread struct {
	ID            int      `db:"id"`
	UserID        int      `db:"user_id"`
	WithUserID    int      `db:"with_user_id"`
	WithUsername  string   `db:"with_username"`
	CreatedAt     NullTime `db:"created_at"`
	LatestMessage string   `db:"latest_message"`
}

func (t *Thread) ToPublic() *PublicThread {
	return &PublicThread{
		WithUsername:  t.WithUsername,
		CreatedAt:     t.CreatedAt.Time,
		LatestMessage: t.LatestMessage,
	}
}

type PublicThread struct {
	WithUsername  string    `json:"with_username"`
	CreatedAt     time.Time `json:"created_at"`
	LatestMessage string    `json:"latest_message"`
}

// NewThread creates a new thread
func NewThread(uid int, withuid int, withuname string, msg string) *Thread {
	return &Thread{
		UserID:       uid,
		WithUserID:   withuid,
		WithUsername: withuname,
		CreatedAt: NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		LatestMessage: msg,
	}
}

// ThreadService defines the protocol for threads
type ThreadService interface {
	GetByUserID(uid int, offset int) ([]*Thread, error)
	Upsert(t *Thread) (int64, error)
}
