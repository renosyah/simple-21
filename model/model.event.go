package model

import "encoding/json"

const (
	// FOR BROADCASTER
	LOBBY_EVENT_ROOM_CREATED = "LOBBY_EVENT_ROOM_CREATED"
	LOBBY_EVENT_ROOM_REMOVE  = "LOBBY_EVENT_ROOM_REMOVED"

	// FOR BROADCASTER RECEIVER
	LOBBY_EVENT_ON_JOIN         = "LOBBY_EVENT_ON_PLAYER_JOIN"
	LOBBY_EVENT_ON_DISCONNECTED = "LOBBY_EVENT_ON_PLAYER_DISCONNECTED"
	LOBBY_EVENT_ON_LOGOUT       = "LOBBY_EVENT_ON_PLAYER_LOGOUT"
	ROOM_EVENT_ON_JOIN          = "ROOM_EVENT_ON_PLAYER_JOIN"
	ROOM_EVENT_ON_DISCONNECTED  = "ROOM_EVENT_ON_PLAYER_DISCONNECTED"
)

type (
	EventData struct {
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}
)

func (_ *EventData) FromJson(b []byte) EventData {
	var e EventData
	if err := json.Unmarshal(b, &e); err != nil {
		return e
	}
	return e
}
