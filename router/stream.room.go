package router

import (
	"context"
	"fmt"
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
		RoomPlayers:  player,
		Cards:        model.NewCards(),
	}

	h.RoomsConn[room.ID] = h.createRoomHub(room)
}

func (h *RouterHub) closeRoom(roomID string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.RoomsConn[roomID]; ok {
		close(h.RoomsConn[roomID].EventBroadcast)
		delete(h.RoomsConn, roomID)
	}
}

func (h *RoomConn) addPlayerRoomConnection(p *RoomPlayerConn) (stream chan model.RoomEventData) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	stream = make(chan model.RoomEventData)
	h.RoomPlayersConn[p.RoomPlayer.ID] = p
	h.RoomPlayersConn[p.RoomPlayer.ID].EventReceiver = stream

	return
}

func (h *RoomConn) removePlayerRoomConnection(p *RoomPlayerConn) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.RoomPlayersConn[p.RoomPlayer.ID]; ok {
		close(h.RoomPlayersConn[p.RoomPlayer.ID].EventReceiver)
		delete(h.RoomPlayersConn, p.RoomPlayer.ID)
	}
}

func (h *RoomConn) receiveBroadcastsEvent(ctx context.Context, wsconn *websocket.Conn, player *RoomPlayerConn) {
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

func (h *RouterHub) createRoomHub(room model.Room) *RoomConn {
	r := &RoomConn{
		Room:            room,
		RoomPlayersConn: make(map[string]*RoomPlayerConn),
		EventBroadcast:  make(chan model.RoomEventData),
	}
	go func() {
		for {
			select {
			case msg := <-r.EventBroadcast:
				switch msg.Status {
				case ROOM_STATUS_USE:

					r.ConnectionMx.RLock()
					for _, c := range r.RoomPlayersConn {
						select {
						case c.EventReceiver <- msg:

						case <-time.After((1 * time.Second)):
							r.removePlayerRoomConnection(c)
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
