package model

import "encoding/json"

const (
	LOBBY_EVENT_JOIN      = "LOBBY_EVENT_PLAYER_JOIN"
	LOBBY_EVENT_NOT_FOUND = "LOBBY_EVENT_PLAYER_NOT_FOUND"
	LOBBY_EVENT_EXIT      = "LOBBY_EVENT_PLAYER_EXIT"
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
