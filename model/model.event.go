package model

import "encoding/json"

const (
	// FOR BROADCASTER
	LOBBY_EVENT_REJOIN = "LOBBY_EVENT_PLAYER_REJOIN"
	LOBBY_EVENT_EXIT   = "LOBBY_EVENT_PLAYER_EXIT"

	LOBBY_EVENT_ROOM_CREATED = "LOBBY_EVENT_ROOM_CREATED"
	ROOM_EVENT_ROOM_REMOVE   = "LOBBY_EVENT_ROOM_CREATED"

	// FOR BROADCASTER RECEIVER
	LOBBY_EVENT_ON_JOIN      = "LOBBY_EVENT_ON_PLAYER_JOIN"
	LOBBY_EVENT_ON_NOT_FOUND = "LOBBY_EVENT_PLAYER_NOT_FOUND"
	LOBBY_EVENT_ON_EXIT      = "LOBBY_EVENT_ON_PLAYER_EXIT"
	LOBBY_EVENT_ON_LOGOUT    = "LOBBY_EVENT_ON_PLAYER_LOGOUT"
)

type (
	EventData struct {
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}
)

func (e *EventData) FromJson(b []byte) {
	if err := json.Unmarshal(b, &e); err != nil {
		e = &EventData{}
	}
}
