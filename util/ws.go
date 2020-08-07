package util

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type hub struct {
	ConnectionMx sync.RWMutex
	Connection   map[*connection]struct{}
	Broadcast    chan []byte
}

type connection struct {
	send chan []byte
	h    *hub
}

type WebsocketServer struct {
	Hub *hub
}

func (w *WebsocketServer) Init() *WebsocketServer {
	w.Hub = NewHub()
	return w
}

func (c *connection) reader(wg *sync.WaitGroup, wsconn *websocket.Conn) {
	defer wg.Done()

	for {
		_, msg, err := wsconn.ReadMessage()
		if err != nil {
			break
		}

		c.h.Broadcast <- msg

	}
}

func (c *connection) write(wg *sync.WaitGroup, wsconn *websocket.Conn) {
	defer wg.Done()

	for msg := range c.send {
		err := wsconn.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			break
		}
	}
}

func NewHub() *hub {
	h := &hub{
		ConnectionMx: sync.RWMutex{},
		Broadcast:    make(chan []byte),
		Connection:   make(map[*connection]struct{}),
	}
	go func() {
		for {
			msg := <-h.Broadcast
			h.ConnectionMx.RLock()
			for c := range h.Connection {
				select {
				case c.send <- msg:

				case <-time.After((1 * time.Second)):
					h.removeconnection(c)
				}
			}
			h.ConnectionMx.RUnlock()
		}

	}()
	return h
}

func (h *hub) addconnection(conn *connection) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	h.Connection[conn] = struct{}{}
}

func (h *hub) removeconnection(conn *connection) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.Connection[conn]; ok {
		delete(h.Connection, conn)
		close(conn.send)
	}
}

func (wsh *WebsocketServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	wsconn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), h: wsh.Hub}
	c.h.addconnection(c)
	defer c.h.removeconnection(c)

	var wg sync.WaitGroup
	wg.Add(2)
	go c.write(&wg, wsconn)
	go c.reader(&wg, wsconn)

	wg.Wait()
	wsconn.Close()

}
