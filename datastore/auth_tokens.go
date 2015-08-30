package datastore

import "gitlab.com/wujiang/asapp"

type authtokenStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(asapp.AuthToken{}, "auth_tokens")
	tm.SetKeys(true, "id")
	tm.ColMap("AccessKeyID").SetUnique(true)
}

func (as *authtokenStore) GetByAccessKeyID(key string) (*asapp.AuthToken, error) {
	var t asapp.AuthToken
	err := as.dbh.SelectOne(&t, `select * from auth_tokens where access_key_id = $1`, key)
	return &t, err
}

func (as *authtokenStore) Create(t *asapp.AuthToken) error {
	err := as.dbh.Insert(t)
	return err
}

func (as *authtokenStore) Update(t *asapp.AuthToken) (int64, error) {
	count, err := as.dbh.Update(t)
	return count, err
}
