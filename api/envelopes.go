package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/wujiang/asapp"
)

func serveGetEnvelopes(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	withUsername := mux.Vars(r)["username"]
	// TODO: page/offset
	page := 1
	if page < 1 {
		return asapp.NewUserError("invalid page number")
	}
	offset := page - 1
	// TODO: get uid from request
	activeUserID := 1
	withUser, err := store.UserStore.GetByUsername(withUsername)
	if err != nil {
		return asapp.NewUserError("unknown username")
	}

	envelopes, err := store.EnvelopeStore.GetByUserIDWithUserID(activeUserID,
		withUser.ID, offset)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	var pubEnv []*asapp.PublicEnvelope
	for _, env := range envelopes {
		pubEnv = append(pubEnv, env.ToPublic())
	}

	return writeJSON(w, asapp.NewJSONResult(pubEnv, page))
}
