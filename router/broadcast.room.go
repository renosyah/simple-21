package router

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
	uuid "github.com/satori/go.uuid"
)

const (
	ROOM_STATUS_USE     = 0
	ROOM_STATUS_NOT_USE = 1
)

func (h *RouterHub) openRoom(pHostID string, player []model.RoomPlayer) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	room := model.Room{
		ID:           fmt.Sprint(uuid.NewV4()),
		PlayerTurnID: pHostID,
		OwnerID:      pHostID,
		RoomPlayers:  player,
		Cards:        model.NewCards(),
	}

	h.Rooms[room.ID] = h.createRoomHub(room)
}

func (h *RouterHub) closeRoom(roomID string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.Rooms[roomID]; ok {
		close(h.Rooms[roomID].EventBroadcast)
		delete(h.Rooms, roomID)
	}
}

func (h *RoomsHub) addPlayerRoomConnection(id string) (stream chan model.RoomEventData) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	stream = make(chan model.RoomEventData)
	h.RoomPlayersConn[id] = stream

	return
}

func (h *RoomsHub) removePlayerRoomConnection(id string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.RoomPlayersConn[id]; ok {
		close(h.RoomPlayersConn[id])
		delete(h.RoomPlayersConn, id)
	}
}

func (h *RoomsHub) receiveBroadcastsEvent(ctx context.Context, wsconn *websocket.Conn, id string) {
	streamClient := h.addPlayerRoomConnection(id)
	defer h.removePlayerRoomConnection(id)

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

func (h *RouterHub) createRoomHub(room model.Room) *RoomsHub {
	r := &RoomsHub{
		ConnectionMx:    sync.RWMutex{},
		Room:            room,
		RoomPlayersConn: make(map[string]chan model.RoomEventData),
		EventBroadcast:  make(chan model.RoomEventData),
	}
	go func() {
		for {
			select {
			case msg := <-r.EventBroadcast:
				switch msg.Status {
				case ROOM_STATUS_USE:

					r.ConnectionMx.RLock()
					for i, c := range r.RoomPlayersConn {
						select {
						case c <- msg:

						case <-time.After((1 * time.Second)):
							r.removePlayerRoomConnection(i)
						}
					}
					r.ConnectionMx.RUnlock()

				case ROOM_STATUS_NOT_USE:

					h.closeRoom(room.ID)
					return

				default:
				}
			default:
			}
		}
	}()
	return r
}
