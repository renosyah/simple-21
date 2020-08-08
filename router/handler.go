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
		TurnPost    int
		TurnsOrder  []string
		Round       int
		Dealer      *model.RoomPlayer
		RoomPlayers map[string]*model.RoomPlayer
		Cards       map[string]*model.Card

		// event in room
		RoomPlayersConn map[string]chan model.RoomEventData
		EventBroadcast  chan model.RoomEventData
	}

	LobbiesHub struct {
		ConnectionMx sync.RWMutex

		// event in lobby
		PlayersConn    map[string]chan model.EventData
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
		Config:       cfg,
		Lobbies:      lobHub,
		Players:      make(map[string]*model.Player),
		Rooms:        make(map[string]*RoomsHub),
	}

	go h.dropOffPlayer()
	go h.dropEmptyRoom()

	return h
}

func (h *RouterHub) dropOffPlayer() {
	for {
		h.ConnectionMx.Lock()
		for _, p := range h.Players {
			if _, ok := h.Lobbies.PlayersConn[p.ID]; !ok {
				delete(h.Players, p.ID)
				h.Lobbies.EventBroadcast <- model.EventData{Name: model.LOBBY_EVENT_ON_LOGOUT}
				break
			}
		}
		h.ConnectionMx.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func (h *RouterHub) dropEmptyRoom() {
	for {
		for id, room := range h.Rooms {
			if len(room.RoomPlayers) == 0 {
				h.Rooms[id].EventBroadcast <- model.RoomEventData{Status: ROOM_STATUS_NOT_USE}
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
}
