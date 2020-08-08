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
	RoomsHub struct {
		ConnectionMx    sync.RWMutex
		Room            model.Room
		RoomPlayersConn map[string]chan model.RoomEventData
		EventBroadcast  chan model.RoomEventData
	}

	LobbiesHub struct {
		ConnectionMx   sync.RWMutex
		PlayersConn    map[string]chan model.EventData
		EventBroadcast chan model.EventData
	}

	RouterHub struct {
		ConnectionMx sync.RWMutex
		Players      map[string]*model.Player
		Lobbies      *LobbiesHub
		Rooms        map[string]*RoomsHub
	}
)

func NewRouterHub() *RouterHub {
	lobHub := &LobbiesHub{
		ConnectionMx:   sync.RWMutex{},
		PlayersConn:    make(map[string]chan model.EventData),
		EventBroadcast: make(chan model.EventData),
	}
	go func() {
		for {
			msg := <-lobHub.EventBroadcast
			lobHub.ConnectionMx.RLock()
			for i, c := range lobHub.PlayersConn {
				select {
				case c <- msg:
				case <-time.After((1 * time.Second)):
					lobHub.removePlayerConnection(i)
				default:
				}

			}
			lobHub.ConnectionMx.RUnlock()
		}

	}()

	h := &RouterHub{
		ConnectionMx: sync.RWMutex{},
		Lobbies:      lobHub,
		Players:      make(map[string]*model.Player),
		Rooms:        make(map[string]*RoomsHub),
	}

	return h
}
