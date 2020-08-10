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
		ConnectionMx sync.RWMutex

		// info room
		Room model.Room

		// data in room
		TurnPost       int
		TurnsOrder     []string
		Round          int
		Status         int
		Dealer         *model.RoomPlayer
		RoomPlayers    map[string]*model.RoomPlayer
		Cards          map[string]*model.Card
		SessionExpired time.Time

		// event in room
		RoomSubscriber map[string]chan model.RoomEventData
		EventBroadcast chan model.RoomEventData
	}

	LobbiesHub struct {
		ConnectionMx sync.RWMutex

		// event in lobby
		Subscriber     map[string]chan model.EventData
		EventBroadcast chan model.EventData
	}

	RouterHub struct {
		ConnectionMx sync.RWMutex
		Players      map[string]*model.Player
		Lobbies      *LobbiesHub
		Rooms        map[string]*RoomsHub
		Config       model.GameConfig
	}
)

func NewRouterHub(cfg model.GameConfig) *RouterHub {
	lobHub := &LobbiesHub{
		ConnectionMx:   sync.RWMutex{},
		Subscriber:     make(map[string]chan model.EventData),
		EventBroadcast: make(chan model.EventData),
	}
	go func() {
		for {
			msg := <-lobHub.EventBroadcast
			lobHub.ConnectionMx.RLock()
			for _, subReceiver := range lobHub.Subscriber {
				select {
				case subReceiver <- msg:
				default:
				}

			}
			lobHub.ConnectionMx.RUnlock()
		}

	}()

	h := &RouterHub{
		ConnectionMx: sync.RWMutex{},
		Config:       cfg,
		Lobbies:      lobHub,
		Players:      make(map[string]*model.Player),
		Rooms:        make(map[string]*RoomsHub),
	}

	go h.dropOffPlayer()
	go h.dropEmptyRoom()

	return h
}
