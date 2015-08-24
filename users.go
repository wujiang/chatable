package asapp

import "time"

const (
	UserClassAdmin = "0"
	UserClassUser  = "1"
)

// User is the corresponding type for a row in users table
type User struct {
	ID            int       `db:"id"`
	Username      string    `db:"username"`
	FirstName     string    `db:"first_name"`
	LastName      string    `db:"last_name"`
	Email         string    `db:"email"`
	PhoneNumber   string    `db:"phone_number"`
	Password      string    `db:"password"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	DeactivatedAt time.Time `db:"deactivated_at"`
	OriginalIP    string    `db:"original_ip"`
	UserClass     string    `db:"user_class"`
}

// NewUser creates a new user
func NewUser(fname, lname, uname, pass, email, phone, ip string) *User {
	return &User{
		FirstName:   fname,
		LastName:    lname,
		Username:    uname,
		Password:    GenerateHash(pass),
		Email:       email,
		PhoneNumber: phone,
		OriginalIP:  ip,
		UserClass:   UserClassUser,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}
}

// UserService defines the protocol for users
type UserService interface {
	GetByUsername(uname string) (*User, error)
	Create(u *User) error
	Update(u *User) (int64, error)
}
