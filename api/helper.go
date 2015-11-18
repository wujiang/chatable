package api

import (
	"encoding/json"
	"net/http"

	"github.com/wujiang/chatable"
)

func writeJSON(w http.ResponseWriter, v interface{}) chatable.CompoundError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return chatable.NewServerError(err.Error())
	} else {
		return nil
	}
}
