package util

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	Host     string
	Port     int
	EndPoint string
}

type WebsocketResponse struct {
	OnConnected   func()
	OnMessage     func(message []byte)
	OnDisconected func()
}

func (c *WebsocketClient) Receiving(ctx context.Context, r *WebsocketResponse) error {

	dialer := websocket.Dialer{
		Subprotocols:    []string{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	url := fmt.Sprintf("ws://%s:%d%s", c.Host, c.Port, c.EndPoint)
	header := http.Header{"Accept-Encoding": []string{"gzip"}}

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}

	r.OnConnected()

	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			break
		default:

			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			r.OnMessage(message)

		}
	}

	r.OnDisconected()

	return nil
}

func (c *WebsocketClient) Send(message []byte) error {

	dialer := websocket.Dialer{
		Subprotocols:    []string{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	url := fmt.Sprintf("ws://%s:%d%s", c.Host, c.Port, c.EndPoint)
	header := http.Header{"Accept-Encoding": []string{"gzip"}}

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return nil

}
