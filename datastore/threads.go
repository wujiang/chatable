package datastore

import "github.com/wujiang/chatable"

type threadStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(chatable.Thread{}, "threads")
	tm.SetKeys(true, "id")
	tm.ColMap("UserID").SetNotNull(true)
	tm.ColMap("WithUserID").SetNotNull(true)
	tm.SetUniqueTogether("user_id", "with_user_id")
}

// Implement the ThreadService

// GetByUserID returns a list of thread for user_id
func (ts *threadStore) GetByUserID(uid int, offset int) ([]*chatable.Thread, error) {
	var t chatable.Thread
	query := `
                select *
                from threads
                where user_id = $1
                order by created_at desc
                limit $2 offset $3
                `
	threads := []*chatable.Thread{}
	ts_, err := ts.dbh.Select(&t, query, uid, chatable.PerPage, offset)
	if err != nil {
		return threads, err
	}
	for _, t_ := range ts_ {
		threads = append(threads, t_.(*chatable.Thread))
	}
	return threads, err
}

// Upsert upserts a row in threads based on (user_id, with_user_id).
// updates doesn't return the rows affected.
func (ts *threadStore) Upsert(t *chatable.Thread) (int64, error) {
	query := `
                with upsert as (
                        update threads
                        set created_at = $1, latest_message = $2
                        where user_id = $3 and with_user_id = $4
                        returning *)
                insert into threads (user_id, with_user_id, author_username,
                        created_at, latest_message)
                        select $3, $4, $5, $1, $2
                where not exists (select * from upsert)
                `
	result, err := ts.dbh.Exec(query, t.CreatedAt, t.LatestMessage,
		t.UserID, t.WithUserID, t.AuthorUsername)
	if err != nil {
		return int64(0), err
	}
	return result.RowsAffected()
}
