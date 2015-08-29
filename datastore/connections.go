package datastore

import "gitlab.com/wujiang/asapp"

type connectionStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(asapp.Connection{}, "connections")
	tm.SetKeys(true, "id")
	tm.ColMap("UserID").SetNotNull(true)
}

// TODO
func (cs *connectionStore) GetByUserID(uid int) ([]*asapp.Connection, error) {
	return []*asapp.Connection{}, nil
}

// TODO
func (cs *connectionStore) Delete(c *asapp.Connection) (int64, error) {
	return cs.dbh.Delete(c)
}
