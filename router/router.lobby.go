package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) HandleLobby(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer wsconn.Close()

	// testing broadcast
	go h.receiveBroadcastsEvent(ctx, wsconn, &PlayerConn{Player: model.Player{ID: "00001", Name: "Reno"}})

	for {

		mType, msg, err := wsconn.ReadMessage()
		if err != nil {
			break
		}

		if mType != websocket.TextMessage {
			break
		}

		event := model.FromJson(msg)
		switch event.Name {
		case model.EVENT_JOIN:
		case model.EVENT_EXIT:
		case model.EVENT_NOT_FOUND:
		default:
		}

		// testing broadcast
		h.EventBroadcast <- model.FromJson(msg)

	}

}
