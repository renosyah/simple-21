package router

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type (
	PlayerConn struct {
		Player        model.Player
		EventReceiver chan model.EventData
	}

	RoomConn struct {
		ConnectionMx   sync.RWMutex
		PlayersConn    map[string]*PlayerConn
		Room           model.Room
		EventBroadcast chan model.EventData
	}

	RouterHub struct {
		ConnectionMx   sync.RWMutex
		PlayersConn    map[string]*PlayerConn
		RoomsConn      map[string]*RoomConn
		EventBroadcast chan model.EventData
	}
)

func NewRouterHub() *RouterHub {
	h := &RouterHub{
		ConnectionMx:   sync.RWMutex{},
		PlayersConn:    make(map[string]*PlayerConn),
		RoomsConn:      make(map[string]*RoomConn),
		EventBroadcast: make(chan model.EventData),
	}
	go func() {
		for {
			msg := <-h.EventBroadcast
			h.ConnectionMx.RLock()
			for i, c := range h.PlayersConn {
				select {
				case c.EventReceiver <- msg:
				case <-time.After((1 * time.Second)):
					h.removePlayerConnection(i)
				}
			}
			h.ConnectionMx.RUnlock()
		}

	}()
	return h
}
