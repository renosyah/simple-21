package model

import "encoding/json"

type (
	Player struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Money int    `json:"money"`
	}
)

func (p *Player) FromJson(b []byte) {
	if err := json.Unmarshal(b, &p); err != nil {
		p = &Player{}
	}
}
