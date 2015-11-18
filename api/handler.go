package api

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/wujiang/chatable"
	"github.com/wujiang/chatable/datastore"
	"github.com/wujiang/chatable/rds"
	"github.com/wujiang/chatable/router"

	goerrors "github.com/go-errors/errors"
)

var (
	store   = datastore.NewDataStore(nil)
	rdsPool = rds.NewRdsPool(nil)
)

func Handler() *mux.Router {
	m := router.API()

	m.Get(router.Register).Handler(handler(serveRegister))
	m.Get(router.CreateAuthToken).Handler(handler(serveCreateAuthToken))
	m.Get(router.DeactivateAuthToken).
		Handler(Authenticate(handler(serveDeactivateAuthToken)))
	m.Get(router.GetInbox).Handler(Authenticate(handler(serveGetThreads)))
	m.Get(router.GetThread).Handler(Authenticate(handler(serveGetEnvelopes)))

	m.Get(router.WSConnect).Handler(Authenticate(handler(serveWSConnect)))

	return m
}

type handler func(http.ResponseWriter, *http.Request) chatable.CompoundError

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err == nil {
		return
	}
	// add stacktrace for errors
	goerr := goerrors.New(err.Error())
	switch err.(type) {
	case chatable.UserError:
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, chatable.NewErrorJSONResult(chatable.JSONError{
			Code:    http.StatusBadRequest,
			Message: "User error",
			Errors:  chatable.ErrorDetails{"error": err.Error()},
		}))
		glog.Warning(goerr.ErrorStack())
	case chatable.FormError:
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, chatable.NewErrorJSONResult(chatable.JSONError{
			Code:    http.StatusBadRequest,
			Message: "Form error",
			Errors:  err.Details(),
		}))
		glog.Warning(goerr.ErrorStack())
	case chatable.AuthenticationError:
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, chatable.NewErrorJSONResult(chatable.JSONError{
			Code:    http.StatusUnauthorized,
			Message: "Authentication error",
			Errors:  chatable.ErrorDetails{"error": err.Error()},
		}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, chatable.NewErrorJSONResult(chatable.JSONError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Errors: chatable.ErrorDetails{
				"error": "internal server error",
			},
		}))
		fmt.Fprint(w, fmt.Sprintf("Internal server error: %s", err))
		glog.Warning(goerr.ErrorStack())
	}
}
