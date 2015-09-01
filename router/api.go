package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()

	m.Path("/register").Methods("POST").Name(Register)
	m.Path("/auth_token").Methods("POST").Name(CreateAuthToken)
	m.Path("/auth_token").Methods("DELETE").Name(DeactivateAuthToken)
	m.Path("/inbox").Methods("GET").Name(GetInbox)
	m.Path("/thread/{username:[a-zA-Z]\\w+}").Methods("GET").
		Name(GetThread)

	m.Path("/ws").Methods("GET").Name(WSConnect)

	return m
}
