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

func (h *RouterHub) openRoom(pHostID, name string, player []model.Player) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	roomPLayers := []model.RoomPlayer{}
	for _, p := range player {
		roomPLayers = append(roomPLayers, model.RoomPlayer{
			ID:       p.ID,
			Name:     p.Name,
			Status:   model.PLAYER_STATUS_INVITED,
			IsOnline: true,
		})
	}

	room := model.Room{
		ID:      fmt.Sprint(uuid.NewV4()),
		OwnerID: pHostID,
		Name:    name,
		Players: roomPLayers,
	}
	hub := h.createRoomHub(room)

	h.Rooms[room.ID] = hub
}

func (h *RouterHub) closeRoom(roomID string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.Rooms[roomID]; ok {
		close(h.Rooms[roomID].EventBroadcast)
		delete(h.Rooms, roomID)
	}
}

func (h *RoomsHub) subscribeRoom(id string) (stream chan model.RoomEventData) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	stream = make(chan model.RoomEventData)
	h.RoomSubscriber[id] = stream

	return
}

func (h *RoomsHub) unSubscribeRoom(id string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.RoomSubscriber[id]; ok {
		close(h.RoomSubscriber[id])
		delete(h.RoomSubscriber, id)
	}
}

func (h *RoomsHub) receiveBroadcastsEvent(ctx context.Context, wsconn *websocket.Conn, id string) {
	subReceiver := h.subscribeRoom(id)
	defer h.unSubscribeRoom(id)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-subReceiver:
			if err := wsconn.WriteMessage(websocket.TextMessage, model.ToJson(msg)); err != nil {
				return
			}
		}
	}
}

func (h *RouterHub) createRoomHub(room model.Room) *RoomsHub {

	timeSet := time.Now().Local()
	timeExp := timeSet.Add(time.Hour*time.Duration(0) +
		time.Minute*time.Duration(h.Config.RoomSessionTime) +
		time.Second*time.Duration(0))

	r := &RoomsHub{
		ConnectionMx: sync.RWMutex{},
		Room:         room,
		TurnPost:     0,
		TurnsOrder:   []string{},
		Dealer: &model.RoomPlayer{
			ID:   fmt.Sprint(uuid.NewV4()),
			Name: fmt.Sprintf("Dealer %s", room.Name),
		},
		RoomPlayers:    make(map[string]*model.RoomPlayer),
		Cards:          make(map[string]*model.Card),
		SessionExpired: timeExp,
		RoomSubscriber: make(map[string]chan model.RoomEventData),
		EventBroadcast: make(chan model.RoomEventData),
	}

	for _, p := range room.Players {
		r.RoomPlayers[p.ID] = &model.RoomPlayer{ID: p.ID, Name: p.Name}
		r.TurnsOrder = append(r.TurnsOrder, p.ID)
	}

	cards := model.NewCards()
	for _, c := range cards {
		r.Cards[c.ID] = &c
	}

	go func() {
		for {
			select {
			case msg := <-r.EventBroadcast:
				switch msg.Status {
				case model.ROOM_STATUS_USE:

					r.ConnectionMx.RLock()
					for _, subReceiver := range r.RoomSubscriber {
						select {
						case subReceiver <- msg:
						default:
						}
					}
					r.ConnectionMx.RUnlock()

				case model.ROOM_STATUS_NOT_USE:

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
