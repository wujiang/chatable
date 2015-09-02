package api

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"gitlab.com/wujiang/asapp"
	"gitlab.com/wujiang/asapp/auth"
)

const (
	BufferedChanLen int = 10
	heartBeat           = 5 * time.Minute
)

var (
	msgRE    = regexp.MustCompile(`^\s*([a-zA-Z]+\w*)\s*:\s*(.*)$`)
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 2048,
	}

	Hub = hub{
		connections: make(map[string][]*connection),
		outgoing:    make(chan asapp.PublicEnvelope),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
)

type connection struct {
	id     int
	conn   *websocket.Conn
	uname  string
	uid    int
	outbuf chan asapp.PublicEnvelope
}

func (c *connection) read() {
	for {
		var env asapp.PublicEnvelope
		// this blocks until new data comes in
		err := c.conn.ReadJSON(&env)
		if err == io.EOF {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		} else if err != nil {
			glog.Warning(err.Error())
			continue
		}
		Hub.outgoing <- asapp.PublicEnvelope{
			Author:      env.Author,
			Recipient:   env.Recipient,
			Message:     env.Message,
			MessageType: asapp.MessageTypeText,
			CreatedAt:   time.Now().UTC(),
		}
	}
}

// close does a cleanup of a connection by closing the outbound channel.
func (c *connection) close() {
	close(c.outbuf)
}

func (c *connection) write() {
	ticker := time.NewTicker(heartBeat)
	defer ticker.Stop()
	defer c.close()

	for {
		select {
		case message, ok := <-c.outbuf:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage,
					[]byte{})
				return
			}
			c.conn.WriteJSON(message)
		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.PingMessage,
				[]byte{})
			if err != nil {
				return
			}
		}
	}
}

type hub struct {
	connections map[string][]*connection
	outgoing    chan asapp.PublicEnvelope
	register    chan *connection
	unregister  chan *connection
}

// Run manages all the channels.
func (h *hub) Run(queue string) {
	for {
		select {
		case c := <-h.register:
			glog.Info(fmt.Sprintf("new connection from %s", c.uname))
			h.connections[c.uname] = append(h.connections[c.uname],
				c)

		case c := <-h.unregister:
			conns := h.connections[c.uname]
			newConns := []*connection{}
			for _, cn := range conns {
				if cn.conn == c.conn {
					c.close()
					glog.Info(fmt.Sprintf("close 1 connection from %s", cn.uname))
				} else {
					newConns = append(newConns, cn)
				}
			}
			h.connections[c.uname] = newConns
		case m := <-h.outgoing:
			if err := rdsPool.Enqueue(queue, m); err != nil {
				glog.Error(err)
			}
		}

	}
}

func serveWSConnect(w http.ResponseWriter, r *http.Request) asapp.CompoundError {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return asapp.NewServerError(err.Error())
	}
	activeUser := auth.ActiveUser(r)
	// this should never happen
	if activeUser == nil {
		return auth.ErrUnauthenticated
	}
	c := &connection{
		conn:   ws,
		uname:  activeUser.Username,
		uid:    activeUser.ID,
		outbuf: make(chan asapp.PublicEnvelope),
	}

	Hub.register <- c
	go c.write()
	c.read()

	return nil
}
