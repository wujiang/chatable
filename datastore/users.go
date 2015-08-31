package datastore

import "gitlab.com/wujiang/asapp"

type userStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(asapp.User{}, "users")
	tm.SetKeys(true, "id")
	tm.ColMap("Username").SetUnique(true).SetNotNull(true)
	tm.ColMap("Email").SetUnique(true).SetNotNull(true)
	tm.ColMap("PhoneNumber").SetUnique(true).SetNotNull(true)
}

// Implement the UserService

// GetByID returns a user with the given user_id
func (us *userStore) GetByID(id int) (*asapp.User, error) {
	var u asapp.User
	err := us.dbh.SelectOne(&u, `select * from users where id = $1`, id)
	return &u, err
}

// GetByUsername returns a user with the given username
func (us *userStore) GetByUsername(uname string) (*asapp.User, error) {
	var u asapp.User
	err := us.dbh.SelectOne(&u, `select * from users where username = $1`,
		uname)
	return &u, err
}

// Create adds a new user into database
func (us *userStore) Create(u *asapp.User) error {
	return us.dbh.Insert(u)
}

// Update updates an existing user in database
func (us *userStore) Update(u *asapp.User) (int64, error) {
	return us.dbh.Update(u)
}
