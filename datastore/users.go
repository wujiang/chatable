package datastore

import (
	"strconv"
	"strings"

	"gitlab.com/wujiang/asapp"
)

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

// GetByIDs returns a list of users with the given ids
func (us *userStore) GetByIDs(ids ...int) ([]*asapp.User, error) {
	var u asapp.User
	var users []*asapp.User
	query := "select * from users where id = any($1::integer[])"
	var sids []string
	for _, id := range ids {
		sids = append(sids, strconv.Itoa(id))
	}
	arg := "{" + strings.Join(sids, ",") + "}"
	us_, err := us.dbh.Select(&u, query, arg)
	if err == nil {
		for _, u_ := range us_ {
			users = append(users, u_.(*asapp.User))
		}
	}
	return users, err
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
