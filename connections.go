package asapp

// Connection represent a row in the database
type Connection struct {
	ID           int    `db:"id"`
	UserID       int    `db:"user_id"`
	MessageQueue string `db:"message_queue"`
}

// NewConnection creates a new connection
func NewConnection(uid int, mq string) *Connection {
	return &Connection{
		UserID:       uid,
		MessageQueue: mq,
	}
}

// ConnectionService defines the protocol for connections
type ConnectionService interface {
	GetByUserID(uid int) ([]*Connection, error)
	Create(c *Connection) error
	DeleteByID(cid int) error
	Delete(c *Connection) (int64, error)
}
