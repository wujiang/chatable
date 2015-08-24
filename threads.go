package asapp

import "time"

// Thread represents a row in the threads table
type Thread struct {
	ID            int       `db:"id"`
	UserID        int       `db:"user_id"`
	WithUserID    int       `db:"with_user_id"`
	WithUsername  string    `db:"with_username"`
	CreatedAt     time.Time `db:"created_at"`
	LatestMessage string    `db:"latest_message"`
}

// NewThread creates a new thread
func NewThread(uid int, withuid int, withuname string, msg string) *Thread {
	return &Thread{
		UserID:        uid,
		WithUserID:    withuid,
		WithUsername:  withuname,
		CreatedAt:     time.Now().UTC(),
		LatestMessage: msg,
	}
}

// ThreadService defines the protocol for threads
type ThreadService interface {
	GetByUserID(uid int, offset int) ([]*Thread, error)
	Upsert(t *Thread) (int64, error)
}
