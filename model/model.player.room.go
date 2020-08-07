package model

type (
	PlayerRoom struct {
		Player Player `json:"player"`
		Bet    int    `json:"bet"`
		Cards  []Card `json:"cards"`
	}
)
