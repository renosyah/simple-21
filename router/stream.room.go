package router

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
	uuid "github.com/satori/go.uuid"
)

func (h *RoomConn) addPlayerRoomConnection(playerConn *PlayerConn) (stream chan model.EventData) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	stream = make(chan model.EventData)
	h.PlayersConn[playerConn.Player.ID] = playerConn
	h.PlayersConn[playerConn.Player.ID].EventReceiver = stream

	return
}

func (h *RoomConn) removePlayerRoomConnection(playerConn *PlayerConn) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.PlayersConn[playerConn.Player.ID]; ok {
		close(h.PlayersConn[playerConn.Player.ID].EventReceiver)
		delete(h.PlayersConn, playerConn.Player.ID)
	}
}

func (h *RoomConn) receiveBroadcastsEvent(ctx context.Context, wsconn *websocket.Conn, player *PlayerConn) {
	streamClient := h.addPlayerRoomConnection(player)
	defer h.removePlayerRoomConnection(player)

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

func (h *RouterHub) CreateRoom(pHostID string, players []model.PlayerRoom) *RoomConn {
	r := &RoomConn{
		Room: model.Room{
			ID:           fmt.Sprint(uuid.NewV4()),
			PlayerTurnID: pHostID,
			Players:      players,
			Cards:        model.NewCards(),
		},
		PlayersConn:    make(map[string]*PlayerConn),
		EventBroadcast: make(chan model.EventData),
	}
	go func() {
		for {
			msg := <-r.EventBroadcast
			r.ConnectionMx.RLock()
			for _, c := range r.PlayersConn {
				select {
				case c.EventReceiver <- msg:
				case <-time.After((1 * time.Second)):
					r.removePlayerRoomConnection(c)
				}
			}
			r.ConnectionMx.RUnlock()
		}

	}()
	return r
}
