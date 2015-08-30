package api

import (
	"net/http"

	"gitlab.com/wujiang/asapp"
)

func serveRegister(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	return serveCreateUser(w, r)
}
