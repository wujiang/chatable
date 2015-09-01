package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/wujiang/asapp"
	"gitlab.com/wujiang/asapp/auth"
)

func serveGetEnvelopes(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
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
		return asapp.NewUserError("unknown username")
	} else if err != nil {
		return asapp.NewServerError(err.Error())
	}
	envelopes, err := store.EnvelopeStore.GetByUserIDWithUserID(
		activeUser.ID, withUser.ID, offset)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	var pubEnv []*asapp.PublicEnvelope
	for _, env := range envelopes {
		pubEnv = append(pubEnv, env.ToPublic(store.UserStore))
	}

	return writeJSON(w, asapp.NewJSONResult(pubEnv, page))
}
