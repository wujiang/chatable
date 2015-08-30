package api

import (
	"net/http"

	"gitlab.com/wujiang/asapp"
)

func serveGetThreads(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	// TODO: get uid from request
	activeUserID := 1
	// TODO: page
	page := 1
	if page < 1 {
		return asapp.NewUserError("invalid page number")
	}
	offset := page - 1
	threads, err := store.ThreadStore.GetByUserID(activeUserID, offset)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	var pubThreads []*asapp.PublicThread
	for th := range threads {
		pubThreads = append(pubThreads, th.ToPublic())
	}
	return writeJSON(w, asapp.NewJSONResult(pubThreads, page))
}
