package router

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/renosyah/simple-21/model"
)

func (h *LobbiesHub) subscribe(id string) (stream chan model.EventData) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()

	stream = make(chan model.EventData)
	h.Subscriber[id] = stream

	return
}

func (h *LobbiesHub) unSubscribe(id string) {
	h.ConnectionMx.Lock()
	defer h.ConnectionMx.Unlock()
	if _, ok := h.Subscriber[id]; ok {
		close(h.Subscriber[id])
		delete(h.Subscriber, id)
	}
}

func (h *LobbiesHub) receiveBroadcastsEvent(ctx context.Context, wsconn *websocket.Conn, id string) {
	subReceiver := h.subscribe(id)
	defer h.unSubscribe(id)

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
