package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestWsClient(t *testing.T) {

	wsc := &WebsocketClient{
		Host:     viper.GetString("ws.host"),
		Port:     viper.GetInt("ws.port"),
		EndPoint: viper.GetString("ws.end_point"),
	}

	ctx := context.Background()

	go wsc.Receiving(ctx, &WebsocketResponse{
		OnConnected: func() {
			t.Logf("connected to websocket service...")
		},
		OnMessage: func(message []byte) {
			t.Logf(fmt.Sprintf("message : %s", string(message)))
		},
		OnDisconected: func() {
			t.Logf("disconnected from websocket service...")
		},
	})
	<- ctx.Done()
}
