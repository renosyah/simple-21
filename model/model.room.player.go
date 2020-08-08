package model

type (
	RoomPlayer struct {
		ID     string `json:"id"`
		Player Player `json:"player"`
		Bet    int    `json:"bet"`
		Cards  []Card `json:"cards"`
	}
)
