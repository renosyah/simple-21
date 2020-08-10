package router

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) HandleStreamRoom(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	pID := r.FormValue("id-player")
	rID := r.FormValue("id-room")

	if pID == "" && rID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, ok := h.Rooms[rID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	player, ok := h.Players[pID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	room.EventBroadcast <- model.RoomEventData{Name: model.ROOM_EVENT_ON_JOIN, Data: *player}

	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.setPlayerOnlineStatus(*player, true)
	room.setPlayerOnlineStatus(*player, true)
	defer h.setPlayerOnlineStatus(*player, false)
	defer room.setPlayerOnlineStatus(*player, false)
	defer wsconn.Close()

	go room.receiveBroadcastsEvent(ctx, wsconn, pID)

	for {

		mType, msg, err := wsconn.ReadMessage()
		if err != nil {
			break
		}

		if mType != websocket.TextMessage {
			break
		}

		event := (&model.RoomEventData{}).FromJson(msg)
		switch event.Name {
		case model.ROOM_EVENT_ON_JOIN:
			/* this event is for client */
		case model.ROOM_EVENT_ON_DISCONNECTED:
			/* this event is for client */
		default:
		}
	}
}
