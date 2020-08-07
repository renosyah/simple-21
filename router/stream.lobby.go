package router

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) addPlayerConnection(id string, playerConn *PlayerConn) (stream chan model.EventData) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	stream = make(chan model.EventData)
	h.PlayersConn[id] = playerConn
	h.PlayersConn[id].EventReceiver = stream

	return
}

func (h *RouterHub) removePlayerConnection(id string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.PlayersConn[id]; ok {
		close(h.PlayersConn[id].EventReceiver)
		delete(h.PlayersConn, id)
	}
}

func (h *RouterHub) receiveBroadcastsEvent(ctx context.Context, wsconn *websocket.Conn, player *PlayerConn) {
	streamClient := h.addPlayerConnection(player.Player.ID, player)
	defer h.removePlayerConnection(player.Player.ID)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-streamClient:
			if err := wsconn.WriteMessage(websocket.TextMessage, model.ToJson(msg)); err != nil {
				return
			}
		}
	}

}

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
