package model

import "encoding/json"

type (
	Player struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Money int    `json:"money"`
	}
)

func PlayerFromJson(b []byte) Player {
	var e Player
	if err := json.Unmarshal(b, &e); err != nil {
		return Player{}
	}
	return e
}
