package api

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"gitlab.com/wujiang/asapp"
	"gitlab.com/wujiang/asapp/datastore"
	"gitlab.com/wujiang/asapp/router"

	goerrors "github.com/go-errors/errors"
)

var (
	store = datastore.NewDataStore(nil)
)

func Handler() *mux.Router {
	m := router.API()
	m.Get(router.CreateAuthToken).Handler(handler(serveCreateAuthToken))
	m.Get(router.DeactivateAuthToken).
		Handler(Authenticate(handler(serveDeactivateAuthToken)))
	m.Get(router.GetInbox).Handler(Authenticate(handler(serveGetThreads)))
	m.Get(router.GetThread).Handler(Authenticate(handler(serveGetEnvelopes)))
	return m
}

type handler func(http.ResponseWriter, *http.Request) asapp.CompoundError

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err == nil {
		return
	}
	// add stacktrace for errors
	goerr := goerrors.New(err.Error())
	switch err.(type) {
	case asapp.UserError:
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, asapp.NewErrorJSONResult(asapp.JSONError{
			Code:    http.StatusBadRequest,
			Message: "User errors",
			Errors:  asapp.ErrorDetails{"error": err.Error()},
		}))
		glog.Warning(goerr.ErrorStack())
	case asapp.FormError:
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, asapp.NewErrorJSONResult(asapp.JSONError{
			Code:    http.StatusBadRequest,
			Message: "Form errors",
			Errors:  err.Details(),
		}))
		glog.Warning(goerr.ErrorStack())
	case asapp.AuthenticationError:
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, asapp.NewErrorJSONResult(asapp.JSONError{
			Code:    http.StatusUnauthorized,
			Message: "Authentication errors",
			Errors:  asapp.ErrorDetails{"error": err.Error()},
		}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, asapp.NewErrorJSONResult(asapp.JSONError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server errors",
			Errors:  asapp.ErrorDetails{"error": err.Error()},
		}))
		fmt.Fprint(w, fmt.Sprintf("Internal server error: %s", err))
		glog.Warning(goerr.ErrorStack())
	}
}
