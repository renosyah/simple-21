package router

import (
	"fmt"

	"github.com/renosyah/simple-21/model"
	uuid "github.com/satori/go.uuid"
)

func (h *RouterHub) addRoomConnection(id string, roomConn *RoomConn) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	h.RoomsConn[id] = roomConn
	h.RoomsConn[id].EventReceiver = make(chan model.EventData)
}

func (h *RouterHub) removeRoomConnection(id string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.RoomsConn[id]; ok {
		close(h.RoomsConn[id].EventReceiver)
		delete(h.RoomsConn, id)
	}
}

func (h *RouterHub) CreateRoom(pHostID string, players []model.PlayerRoom) *RoomConn {
	return &RoomConn{
		Room: model.Room{
			ID:           fmt.Sprint(uuid.NewV4()),
			PlayerTurnID: pHostID,
			Players:      players,
			Cards:        model.NewCards(),
		},
	}
}
