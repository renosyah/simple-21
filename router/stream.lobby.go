package router

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) HandleStreamLobby(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	pID := r.FormValue("id-player")

	if pID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	player, ok := h.Players[pID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.setPlayerOnlineStatus(*player, true, true)
	defer h.setPlayerOnlineStatus(*player, false, true)
	defer wsconn.Close()

	go h.Lobbies.receiveBroadcastsEvent(ctx, wsconn, pID)

	for {

		mType, msg, err := wsconn.ReadMessage()
		if err != nil {
			break
		}

		if mType != websocket.TextMessage {
			break
		}

		event := (&model.EventData{}).FromJson(msg)
		switch event.Name {
		case model.LOBBY_EVENT_ON_JOIN:
			/* this event is for client */
		case model.LOBBY_EVENT_ON_DISCONNECTED:
			/* this event is for client */
		case model.LOBBY_EVENT_ON_ROOM_CREATED:
			/* this event is for client */
		case model.LOBBY_EVENT_ON_ROOM_REMOVE:
			/* this event is for client */
		case model.LOBBY_EVENT_ON_LOGOUT:
			/* this event is for client */
		default:
		}
	}
}
