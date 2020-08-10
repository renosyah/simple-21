package model

type (
	RoomBet struct {
		PlayerID string `json:"player_id"`
		RoomID   string `json:"room_id"`
		Bet      int    `json:"bet"`
	}
)
