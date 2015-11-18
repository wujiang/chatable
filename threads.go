package chatable

import "time"

// Thread represents a row in the threads table
type Thread struct {
	ID             int      `db:"id"`
	UserID         int      `db:"user_id"`
	WithUserID     int      `db:"with_user_id"`
	AuthorUsername string   `db:"author_username"`
	CreatedAt      NullTime `db:"created_at"`
	LatestMessage  string   `db:"latest_message"`
}

func (t *Thread) ToPublic() *PublicThread {
	return &PublicThread{
		AuthorUsername: t.AuthorUsername,
		CreatedAt:      t.CreatedAt.Time,
		LatestMessage:  t.LatestMessage,
	}
}

type PublicThread struct {
	AuthorUsername string    `json:"author_username"`
	CreatedAt      time.Time `json:"created_at"`
	LatestMessage  string    `json:"latest_message"`
}

// NewThread creates 2 new threads
func NewThread(uid int, withuid int, author string, msg string) (*Thread, *Thread) {
	dt := NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	return &Thread{
			UserID:         uid,
			WithUserID:     withuid,
			AuthorUsername: author,
			CreatedAt:      dt,
			LatestMessage:  msg,
		}, &Thread{
			UserID:         withuid,
			WithUserID:     uid,
			AuthorUsername: author,
			CreatedAt:      dt,
			LatestMessage:  msg,
		}
}

// ThreadService defines the protocol for threads
type ThreadService interface {
	GetByUserID(uid int, offset int) ([]*Thread, error)
	Upsert(t *Thread) (int64, error)
}
