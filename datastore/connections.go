package datastore

import "gitlab.com/wujiang/asapp"

type connectionStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(asapp.Connection{}, "connections")
	tm.SetKeys(true, "id")
	tm.ColMap("UserID").SetNotNull(true)
	tm.ColMap("MessageQueue").SetNotNull(true)
	tm.SetUniqueTogether("user_id", "message_queue")
}

// GetByUserID returns all connections a user has.
func (cs *connectionStore) GetByUserID(uid int) ([]*asapp.Connection, error) {
	var c asapp.Connection
	connections := []*asapp.Connection{}
	conns, err := cs.dbh.Select(&c,
		"select * from connections where user_id = $1", uid)
	if err != nil {
		return connections, err
	}
	for _, c_ := range conns {
		connections = append(connections, c_.(*asapp.Connection))
	}
	return connections, nil
}

// Create adds a new connection into database
func (cs *connectionStore) Create(c *asapp.Connection) error {
	return cs.dbh.Insert(c)
}

func (cs *connectionStore) DeleteByID(cid int) error {
	_, err := cs.dbh.Exec("delete from connections where id = $1", cid)
	return err
}

// Delete removes a connection row from database
func (cs *connectionStore) Delete(c *asapp.Connection) (int64, error) {
	return cs.dbh.Delete(c)
}
