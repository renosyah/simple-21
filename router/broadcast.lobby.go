package router

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

func (h *LobbiesHub) addPlayerConnection(id string) (stream chan model.EventData) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	stream = make(chan model.EventData)
	h.PlayersConn[id] = stream

	return
}

func (h *LobbiesHub) removePlayerConnection(id string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.PlayersConn[id]; ok {
		close(h.PlayersConn[id])
		delete(h.PlayersConn, id)
	}
}

func (h *RouterHub) setPlayerOnlineStatus(id string, isOnline bool) {
	h.ConnectionMx.Lock()
	if p, ok := h.Players[id]; ok {
		p.IsOnline = isOnline
	}
	h.ConnectionMx.Unlock()

	h.Lobbies.EventBroadcast <- model.EventData{Name: model.LOBBY_EVENT_ON_DISCONNECTED}
}

func (h *LobbiesHub) receiveBroadcastsEvent(ctx context.Context, wsconn *websocket.Conn, id string) {
	streamClient := h.addPlayerConnection(id)
	defer h.removePlayerConnection(id)

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
