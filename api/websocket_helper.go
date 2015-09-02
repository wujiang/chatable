package api

import (
	"github.com/golang/glog"
	"gitlab.com/wujiang/asapp"
)

// just a place holder
type QueueManager struct {
}

// dispatch pops an envelope from shared queue and pushes to
// all servers' queues. queue is the shared queue
func (qm *QueueManager) Dispatch(queue string, key string) {
	for {
		// this is blocking operation
		env, err := rdsPool.Dequeue(queue)
		if err != nil {
			glog.Error(err.Error())
		}

		// persist to db
		if err = asapp.PersistEnvelope(env, store.UserStore,
			store.EnvelopeStore, store.ThreadStore); err != nil {
			glog.Error(err.Error())
		}

		// push to all message queues
		queues, err := rdsPool.QMMembers(key)
		for _, q := range queues {
			err = rdsPool.Enqueue(q, env)
			if err != nil {
				glog.Error(err.Error())
			}
		}
	}
}

func (qm *QueueManager) Pop(msgQueue string) {
	for {
		// this is blocking operation
		env, err := rdsPool.Dequeue(msgQueue)
		if err != nil {
			glog.Error(err.Error())
		}

		conns := Hub.connections[env.Recipient]
		if len(conns) == 0 {
			continue
		}
		newConns := []*connection{}
		for _, c := range conns {
			select {
			case c.outbuf <- env:
				newConns = append(newConns, c)
			default:
				close(c.outbuf)
			}
		}
		Hub.connections[env.Recipient] = newConns

	}
}
