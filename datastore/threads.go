package datastore

import (
	"database/sql"

	"gitlab.com/wujiang/asapp"
)

type threadStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(asapp.Thread{}, "threads")
	tm.SetKeys(true, "id")
	tm.ColMap("UserID").SetNotNull(true)
	tm.ColMap("WithUserID").SetNotNull(true)
	tm.SetUniqueTogether("user_id", "with_user_id")
}

// Implement the ThreadService

// GetByUserID returns a list of thread for user_id
func (ts *threadStore) GetByUserID(uid int, offset int) ([]*asapp.Thread, error) {
	var t asapp.Thread
	query := `
                select *
                from threads
                where user_id = $1
                order by created_at desc
                limit $2 offset $3
                `
	threads := []*asapp.Thread{}
	ts_, err := ts.dbh.Select(&t, query, uid, asapp.PerPage, offset)
	if err != nil {
		return threads, err
	}
	for _, t_ := range ts_ {
		threads = append(threads, t_.(*asapp.Thread))
	}
	return threads, err
}

func (ts *threadStore) getByUserIDWithUserID(uid int, withuid int) (*asapp.Thread, error) {
	var t asapp.Thread
	err := ts.dbh.SelectOne(&t,
		`select 1 from threads where user_id = $1 and with_user_id = $2`,
		uid, withuid)
	return &t, err
}

// Upsert updates or inserts a thread
func (ts *threadStore) Upsert(t *asapp.Thread) (int64, error) {
	_, err := ts.getByUserIDWithUserID(t.UserID, t.WithUserID)
	if err == sql.ErrNoRows {
		return 1, ts.dbh.Insert(t)
	} else if err == nil {
		return ts.dbh.Update(t)
	} else {
		return 0, err
	}
}
