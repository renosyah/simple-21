package router

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
	"github.com/renosyah/simple-21/util"
	uuid "github.com/satori/go.uuid"
)

func (h *RouterHub) openRoom(add model.AddRoom) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	roomPlayers := []model.RoomPlayer{}
	for _, p := range add.Players {
		roomPlayers = append(roomPlayers, model.RoomPlayer{
			ID:       p.ID,
			Name:     p.Name,
			Status:   model.PLAYER_STATUS_INVITED,
			Cards:    []model.Card{},
			Money:    0,
			IsOnline: true,
			IsBot:    false,
		})
	}

	for i := 0; i < add.Bot; i++ {
		roomPlayers = append(roomPlayers, model.RoomPlayer{
			ID:       fmt.Sprint(uuid.NewV4()),
			Name:     fmt.Sprintf("%s (Bot)", util.GenerateRandomName(true)),
			Status:   model.PLAYER_STATUS_INVITED,
			Cards:    []model.Card{},
			Money:    0,
			IsOnline: true,
			IsBot:    true,
		})
	}

	room := model.Room{
		ID:          fmt.Sprint(uuid.NewV4()),
		OwnerID:     add.HostID,
		Name:        add.Name,
		Players:     roomPlayers,
		CardGroups:  add.CardGroups,
		CanDrawCard: true,
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

	cards := model.NewCards(room.CardGroups)

	timeSet := time.Now().Local()
	timeExp := timeSet.Add(time.Hour*time.Duration(0) +
		time.Minute*time.Duration(h.Config.RoomSessionTime) +
		time.Second*time.Duration(0))

	r := &RoomsHub{
		ConnectionMx: sync.RWMutex{},
		Room:         room,
		Turn: &TurnHandler{
			TurnPost:   0,
			TurnsOrder: []string{},
		},
		Dealer: &model.RoomPlayer{
			ID:    fmt.Sprint(uuid.NewV4()),
			Name:  util.GenerateRandomName(true),
			Cards: []model.Card{},
		},
		Round:          1,
		RoomPlayers:    make(map[string]*model.RoomPlayer),
		Cards:          make(map[string]*model.Card),
		Scores:         make(map[int]*model.Score),
		SessionExpired: timeExp,
		RoomSubscriber: make(map[string]chan model.RoomEventData),
		EventBroadcast: make(chan model.RoomEventData),
	}
	for i, p := range room.Players {
		player := &model.RoomPlayer{
			ID:        p.ID,
			Name:      p.Name,
			Cards:     []model.Card{},
			TurnOrder: i,
			Money:     p.Money,
			IsBot:     p.IsBot,
		}

		if p.IsBot {
			r.runBotFunction(h, player)
		}

		r.RoomPlayers[player.ID] = player
		r.Turn.TurnsOrder = append(r.Turn.TurnsOrder, p.ID)
	}

	for _, c := range cards {
		r.Cards[c.ID] = c.CopyPointer()
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
