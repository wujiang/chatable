package datastore

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"gitlab.com/wujiang/asapp"
)

var dbm = &gorp.DbMap{
	Dialect: gorp.PostgresDialect{},
}
var dbHandler gorp.SqlExecutor = dbm

// Init opens a connection to the database
func Init(pg string) {
	if dbm.Db != nil {
		if err := dbm.Db.Ping(); err == nil {
			return
		}
	}
	db, err := sql.Open("postgres", pg)
	if err != nil {
		glog.Fatal(err)
	}
	dbm.Db = db
	if err = dbm.Db.Ping(); err != nil {
		glog.Fatal("Unable to establish a connection to DB")
	}
}

// Exit closes the database connection
func Exit() {
	dbm.Db.Close()
}

// Database is the portal to database
type DataStore struct {
	UserStore       asapp.UserService
	ThreadStore     asapp.ThreadService
	EnvelopeStore   asapp.EnvelopeService
	ConnectionStore asapp.ConnectionService

	dbh gorp.SqlExecutor
}

// NewDataStore creates new datastore
func NewDataStore(dbh gorp.SqlExecutor) *DataStore {
	if dbh == nil {
		dbh = dbHandler
	}
	ds := &DataStore{dbh: dbh}
	ds.UserStore = &userStore{ds}
	ds.ThreadStore = &threadStore{ds}
	ds.EnvelopeStore = &envelopeStore{ds}
	ds.ConnectionStore = &connectionStore{ds}
	return ds
}
