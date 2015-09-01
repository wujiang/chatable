package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"gitlab.com/wujiang/asapp"
	"gitlab.com/wujiang/asapp/auth"
)

const (
	BufferedChanLen int = 10
)

var (
	msgRE    = regexp.MustCompile(`^\s*([a-zA-Z]+\w*)\s*:\s*(.*)$`)
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 2048,
	}

	h = hub{
		connections: make(map[string][]*connection),
		outgoing:    make(chan message),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
)

// parseText parses a text to get recipient and message
// recipient: how are you? --> recipient, how are you?
func parseText(msg string) (string, string) {
	parts := msgRE.FindStringSubmatch(msg)
	if len(parts) != 3 {
		return "", ""
	}
	return parts[1], parts[2]
}

type connection struct {
	conn   *websocket.Conn
	uname  string
	outbuf chan message
}

func (c *connection) read() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		recipient, m := parseText(string(msg))
		fmt.Println("recipient: ", recipient, ", msg: ", m)
		if recipient == "" {
			continue
		}
		h.outgoing <- message{
			sender:    c.uname,
			recipient: recipient,
			msg:       m,
		}
	}
}

func (c *connection) write() {
	for {
		select {
		case message, ok := <-c.outbuf:
			if !ok {
				continue
			}
			msg := fmt.Sprintf("[%s] %s", message.sender, message.msg)
			c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

type message struct {
	sender    string
	recipient string
	msg       string // TODO: other types
}

type hub struct {
	connections map[string][]*connection
	outgoing    chan message
	register    chan *connection
	unregister  chan *connection
}

func (h *hub) exec() {
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
					close(cn.outbuf)
					glog.Info(fmt.Sprintf("close 1 connection from %s", cn.uname))
				} else {
					newConns = append(newConns, cn)
				}
			}
			h.connections[c.uname] = newConns
		case m := <-h.outgoing:
			conns := h.connections[m.recipient]
			newConns := []*connection{}
			for _, c := range conns {
				select {
				case c.outbuf <- m:
					newConns = append(newConns, c)
				default:
					close(c.outbuf)
				}
			}
			h.connections[m.recipient] = newConns
		}

	}
}

func init() {
	go h.exec()
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
		outbuf: make(chan message),
	}

	h.register <- c
	go c.write()
	c.read()

	return nil
}
