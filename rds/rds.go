package rds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"gitlab.com/wujiang/asapp"
)

var (
	// shared pool
	rdsPool = &RdsPool{
		pool: &redis.Pool{},
	}
)

func Init(host string) {
	p := redis.Pool{
		MaxIdle: 5,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host)
			if err != nil {
				glog.Error(err)
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	*rdsPool.pool = p
}

func Exit() {
	rdsPool.pool.Close()
}

type RdsPool struct {
	pool *redis.Pool
}

// NewRdsPool returns a new instance.
func NewRdsPool(pool *redis.Pool) *RdsPool {
	if pool == nil {
		pool = rdsPool.pool
	}
	return &RdsPool{
		pool: pool,
	}
}

// Implementation of RdsService

// Enqueue pushes a PublicEnvelope into the tail of a given queue.
func (r *RdsPool) Enqueue(queue string, env asapp.PublicEnvelope) asapp.CompoundError {
	conn := r.pool.Get()
	defer conn.Close()

	bt, err := json.Marshal(env)
	if err != nil {
		return asapp.NewServerError(fmt.Sprintf("Can not marshal %+v", env))
	}
	if _, err = conn.Do("RPUSH", queue, bt); err != nil {
		return asapp.NewServerError(err.Error())
	}
	return nil
}

// Dequeue pops the first element from the given queue. This is a blocking
// operation.
func (r *RdsPool) Dequeue(queue string) (asapp.PublicEnvelope, asapp.CompoundError) {
	conn := r.pool.Get()
	defer conn.Close()

	var env asapp.PublicEnvelope
	val, err := redis.Values(conn.Do("BLPOP", queue, 0))
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

func (r *RdsPool) AddToQM(key string, queue string) asapp.CompoundError {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, queue)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	return nil
}

func (r *RdsPool) QMMembers(key string) ([]string, asapp.CompoundError) {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return []string{}, asapp.NewServerError(err.Error())
	}
	return val, nil
}

func (r *RdsPool) RemoveFromQM(key string, queue string) asapp.CompoundError {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SREM", key, queue)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	return nil
}
