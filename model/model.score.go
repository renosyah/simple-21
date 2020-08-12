package model

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type (
	Score struct {
		ID      string        `json:"id"`
		Dealer  ScorePlayer   `json:"dealer"`
		Players []ScorePlayer `json:"players"`
		Round   int           `json:"round"`
	}
	ScorePlayer struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Bet    int    `json:"bet"`
		Total  int    `json:"total"`
		Status int    `json:"status"`
	}
)

func (s *Score) Copy() Score {
	return Score{
		ID:      s.ID,
		Dealer:  s.Dealer,
		Players: s.Players,
		Round:   s.Round,
	}
}

func RecordScore(round int, dealer *RoomPlayer, players map[string]*RoomPlayer) *Score {
	plrs := []ScorePlayer{}
	for _, p := range players {
		plrs = append(plrs, ScorePlayer{
			ID:     p.ID,
			Name:   p.Name,
			Bet:    p.Bet,
			Total:  p.Total,
			Status: p.Status,
		})
	}

	return &Score{
		ID: fmt.Sprint(uuid.NewV4()),
		Dealer: ScorePlayer{
			ID:    dealer.ID,
			Name:  dealer.Name,
			Total: dealer.Total,
		},
		Players: plrs,
		Round:   round,
	}
}
