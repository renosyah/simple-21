package model

import "encoding/json"

const (
	EVENT_JOIN      = "EVENT_PLAYER_JOIN"
	EVENT_NOT_FOUND = "EVENT_PLAYER_NOT_FOUND"
	EVENT_EXIT      = "EVENT_PLAYER_EXIT"
)

type (
	EventData struct {
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}
)

func FromJson(b []byte) EventData {
	var e EventData
	if err := json.Unmarshal(b, &e); err != nil {
		return EventData{}
	}
	return e
}
