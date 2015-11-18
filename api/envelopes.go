package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/auth"
)

func serveGetEnvelopes(w http.ResponseWriter, r *http.Request) chatable.CompoundError {
	withUsername := mux.Vars(r)["username"]
	params := r.URL.Query()
	pg := params.Get("page")
	page, err := strconv.Atoi(pg)
	if err != nil || page < 1 {
		page = 1
	}

	offset := page - 1
	activeUser := auth.ActiveUser(r)
	// this should never happen
	if activeUser == nil {
		return auth.ErrUnauthenticated
	}
	withUser, err := store.UserStore.GetByUsername(withUsername)
	if err == sql.ErrNoRows {
		return chatable.NewUserError("unknown username")
	} else if err != nil {
		return chatable.NewServerError(err.Error())
	}
	envelopes, err := store.EnvelopeStore.GetByUserIDWithUserID(
		activeUser.ID, withUser.ID, offset)
	if err != nil {
		return chatable.NewServerError(err.Error())
	}
	var pubEnv []*chatable.PublicEnvelope
	for _, env := range envelopes {
		pubEnv = append(pubEnv, env.ToPublic(store.UserStore))
	}

	return writeJSON(w, chatable.NewJSONResult(pubEnv, page))
}
