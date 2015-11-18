package datastore

import "github.com/wujiang/chatable"

type authtokenStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(chatable.AuthToken{}, "auth_tokens")
	tm.SetKeys(true, "id")
	tm.ColMap("AccessKeyID").SetUnique(true)
}

func (as *authtokenStore) GetByAccessKeyID(key string) (*chatable.AuthToken, error) {
	var t chatable.AuthToken
	query := `select *
                from auth_tokens
                where access_key_id = $1
                        and is_active is true
                        and expires_at >= now()`
	err := as.dbh.SelectOne(&t, query, key)
	return &t, err
}

func (as *authtokenStore) Create(t *chatable.AuthToken) error {
	err := as.dbh.Insert(t)
	return err
}

func (as *authtokenStore) Update(t *chatable.AuthToken) (int64, error) {
	count, err := as.dbh.Update(t)
	return count, err
}
