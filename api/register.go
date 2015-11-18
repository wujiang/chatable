package api

import (
	"net/http"

	"github.com/wujiang/chatable"
)

func serveRegister(w http.ResponseWriter, r *http.Request) chatable.CompoundError {
	return serveCreateUser(w, r)
}
