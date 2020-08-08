package model

import "encoding/json"

type (
	RoomEventData struct {
		Name   string      `json:"name"`
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
	}
)

func (r *RoomEventData) FromJson(b []byte) RoomEventData {
	var e RoomEventData
	if err := json.Unmarshal(b, &e); err != nil {
		return RoomEventData{}
	}
	return e
}
