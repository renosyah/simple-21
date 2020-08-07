package router

import (
	"net/http"
	"sync"

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
		Room          model.Room
		EventReceiver chan model.EventData
	}

	RouterHub struct {
		ConnectionMx   sync.RWMutex
		PlayersConn    map[string]*PlayerConn
		RoomsConn      map[string]*RoomConn
		EventBroadcast chan model.EventData
	}
)
