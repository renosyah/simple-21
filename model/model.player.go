package model

import (
	"encoding/json"
	"time"
)

type (
	Player struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		Money          int       `json:"money"`
		IsOnline       bool      `json:"is_online"`
		SessionExpired time.Time `json:"-"`
	}
)

func (p *Player) FromJson(b []byte) {
	if err := json.Unmarshal(b, &p); err != nil {
		p = &Player{}
	}
}
