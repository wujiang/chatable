package api

import (
	"encoding/json"
	"net/http"

	"gitlab.com/wujiang/asapp"
)

func writeJSON(w http.ResponseWriter, v interface{}) asapp.CompoundError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return asapp.NewServerError(err.Error())
	} else {
		return nil
	}
}
