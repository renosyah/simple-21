package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

func (h *RouterHub) HandleStreamLobby(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer wsconn.Close()

	for {

		mType, msg, err := wsconn.ReadMessage()
		if err != nil {
			return
		}

		if mType != websocket.TextMessage {
			return
		}

		event := &model.EventData{}
		event.FromJson(msg)

		switch event.Name {
		case model.LOBBY_EVENT_REJOIN:

			/* */
			/* for player already register */
			/* and decide to join lobby */
			/* check his data first, if exist */
			/* he have right to receive event */
			/* if not, just return message his data */
			/* not exist */
			/* */

			p := &model.Player{}
			p.FromJson(model.ToJson(event.Data))

			if _, ok := h.Players[p.ID]; ok {

				go h.Lobbies.receiveBroadcastsEvent(ctx, wsconn, p.ID)

				h.Lobbies.EventBroadcast <- model.EventData{
					Name: model.LOBBY_EVENT_ON_JOIN,
					Data: p,
				}

			} else {

				resp := model.EventData{
					Name: model.LOBBY_EVENT_ON_NOT_FOUND,
					Data: model.Player{},
				}
				wsconn.WriteMessage(websocket.TextMessage, model.ToJson(resp))

			}

		case model.LOBBY_EVENT_EXIT:

			/* */
			/* for player already register */
			/* and decide to exit from lobby */
			/* broadcast to other this player */
			/* is decide to disconect */
			/* */

			p := &model.Player{}
			p.FromJson(model.ToJson(event.Data))

			h.Lobbies.EventBroadcast <- model.EventData{
				Name: model.LOBBY_EVENT_ON_EXIT,
				Data: p,
			}

		case model.LOBBY_EVENT_ON_NOT_FOUND:
			/* this event is for client */
		case model.LOBBY_EVENT_ON_JOIN:
			/* this event is for client */
		case model.LOBBY_EVENT_ON_EXIT:
			/* this event is for client */
		default:
		}
	}

}
