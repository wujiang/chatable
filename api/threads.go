package api

import (
	"net/http"
	"strconv"

	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/auth"
)

func serveGetThreads(w http.ResponseWriter, r *http.Request) chatable.CompoundError {
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

	threads, err := store.ThreadStore.GetByUserID(activeUser.ID, offset)
	if err != nil {
		return chatable.NewServerError(err.Error())
	}
	var pubThreads []*chatable.PublicThread
	for _, th := range threads {
		pubThreads = append(pubThreads, th.ToPublic())
	}
	return writeJSON(w, chatable.NewJSONResult(pubThreads, page))
}
