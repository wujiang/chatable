package datastore

import (
	"time"

	"gitlab.com/wujiang/asapp"
)

type envelopeStore struct{ *DataStore }

func init() {
	tm := dbm.AddTableWithName(asapp.Envelope{}, "envelopes")
	tm.SetKeys(true, "id")
	tm.ColMap("UserID").SetNotNull(true)
	tm.ColMap("WithUserID").SetNotNull(true)
	tm.SetUniqueTogether("user_id", "with_user_id")
}

// Implement the EnvelopeService

// GetByUserIDWithUserID returns a list of most recent envelopes
// between user_id and with_user_id.
func (es *envelopeStore) GetByUserIDWithUserID(uid int, withuid int,
	offset int) ([]*asapp.Envelope, error) {
	var env asapp.Envelope
	query := `
                select id, user_id, with_user_id, is_incoming, created_at,
                        deleted_at, read_at, message, message_type
                from envelopes
                where user_id = $1 and with_user_id = $2 and deleted_at is null
                order by created_at desc
                limit $3 offset $4
                `
	envelopes := []*asapp.Envelope{}
	envs, err := es.dbh.Select(&env, query, uid, withuid, asapp.PerPage,
		offset)
	if err != nil {
		return envelopes, err
	}
	for _, e := range envs {
		envelopes = append(envelopes, e.(*asapp.Envelope))
	}
	return envelopes, err
}

// Create adds a new row in envelopes
func (es *envelopeStore) Create(env *asapp.Envelope) error {
	return es.dbh.Insert(env)
}

// MarkDelete marks an envelop as deleted
func (es *envelopeStore) MarkDelete(env *asapp.Envelope) (int64, error) {
	env.DeletedAt = asapp.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	return es.dbh.Update(env)
}

// MarkRead marks an envelop as read
func (es *envelopeStore) MarkRead(env *asapp.Envelope) (int64, error) {
	env.ReadAt = asapp.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	return es.dbh.Update(env)
}
