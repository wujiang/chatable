package rds

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"gitlab.com/wujiang/asapp"
)

var (
	// sharable connection to redis
	conn    redis.Conn
	rdsConn = &RdsConn{
		conn: &conn,
	}
)

func Init(host string) {
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		glog.Fatal(err)
	}
	*rdsConn.conn = conn
}

func Exit() {
	(*rdsConn.conn).Close()
}

type RdsConn struct {
	conn *redis.Conn
}

// NewRdsConn returns a new instance.
func NewRdsConn(conn *redis.Conn) *RdsConn {
	if conn == nil {
		conn = rdsConn.conn
	}
	return &RdsConn{
		conn: conn,
	}
}

// Implementation of RdsService

// Enqueue pushes a PublicEnvelope into the tail of a given queue.
func (r *RdsConn) Enqueue(queue string, env asapp.PublicEnvelope) asapp.CompoundError {
	bt, err := json.Marshal(env)
	if err != nil {
		return asapp.NewServerError(fmt.Sprintf("Can not marshal %+v", env))
	}
	if _, err = (*r.conn).Do("RPUSH", queue, bt); err != nil {
		return asapp.NewServerError(err.Error())
	}
	return nil
}

// Dequeue pops the first element from the given queue. This is a blocking
// operation.
func (r *RdsConn) Dequeue(queue string) (asapp.PublicEnvelope, asapp.CompoundError) {
	var env asapp.PublicEnvelope
	val, err := redis.Values((*r.conn).Do("BLPOP", queue, 0))
	if err != nil {
		return env, asapp.NewServerError(err.Error())
	}
	var q, bt []byte
	if _, err = redis.Scan(val, &q, &bt); err != nil {
		return env, asapp.NewServerError(err.Error())
	}
	if err = json.Unmarshal(bt, &env); err != nil {
		return env, asapp.NewServerError(err.Error())
	}
	return env, nil
}

func (r *RdsConn) AddToQM(key string, queue string) asapp.CompoundError {
	_, err := (*r.conn).Do("SADD", key, queue)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	return nil
}

func (r *RdsConn) QMMembers(key string) ([]string, asapp.CompoundError) {
	val, err := redis.Strings((*r.conn).Do("SMEMBERS", key))
	if err != nil {
		return []string{}, asapp.NewServerError(err.Error())
	}
	return val, nil
}

func (r *RdsConn) RemoveFromQM(key string, queue string) asapp.CompoundError {
	_, err := (*r.conn).Do("SREM", key, queue)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	return nil
}
