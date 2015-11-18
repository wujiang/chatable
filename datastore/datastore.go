package datastore

import (
	"database/sql"

	"github.com/wujiang/chatable"

	"github.com/coopernurse/gorp"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
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
	UserStore      chatable.UserService
	ThreadStore    chatable.ThreadService
	EnvelopeStore  chatable.EnvelopeService
	AuthTokenStore chatable.AuthTokenService

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
	ds.AuthTokenStore = &authtokenStore{ds}
	return ds
}

// CreateTables creates all registered tables. This is used by tests.
func CreateTables() error {
	return dbm.CreateTablesIfNotExists()
}

// DropTables drops all registered tables. This is used by tests.
func DropTables() error {
	return dbm.DropTablesIfExists()
}
