package model

import "encoding/json"

const (

	// FOR BROADCASTER RECEIVER
	LOBBY_EVENT_ON_JOIN         = "LOBBY_EVENT_ON_PLAYER_JOIN"
	LOBBY_EVENT_ON_DISCONNECTED = "LOBBY_EVENT_ON_PLAYER_DISCONNECTED"
	LOBBY_EVENT_ON_LOGOUT       = "LOBBY_EVENT_ON_PLAYER_LOGOUT"
	LOBBY_EVENT_ON_ROOM_CREATED = "LOBBY_EVENT_ON_ROOM_CREATED"
	LOBBY_EVENT_ON_ROOM_REMOVE  = "LOBBY_EVENT_ON_ROOM_REMOVED"
)

type EventData struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func (_ *EventData) FromJson(b []byte) EventData {
	var e EventData
	if err := json.Unmarshal(b, &e); err != nil {
		return e
	}
	return e
}

const (
	// FOR BROADCASTER RECEIVER
	ROOM_EVENT_ON_JOIN           = "ROOM_EVENT_ON_PLAYER_JOIN"
	ROOM_EVENT_ON_DISCONNECTED   = "ROOM_EVENT_ON_PLAYER_DISCONNECTED"
	ROOM_EVENT_ON_PLAYER_SET_BET = "ROOM_EVENT_ON_PLAYER_SET_BET"
	ROOM_EVENT_ON_GAME_START     = "ROOM_EVENT_ON_GAME_START"
	ROOM_EVENT_ON_CARD_GIVEN     = "ROOM_EVENT_ON_CARD_GIVEN"
)

type RoomEventData struct {
	Name   string      `json:"name"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (_ *RoomEventData) FromJson(b []byte) RoomEventData {
	var e RoomEventData
	if err := json.Unmarshal(b, &e); err != nil {
		return RoomEventData{}
	}
	return e
}
