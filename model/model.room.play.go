package model

const (
	ROOM_TURN_CHOOSE_HIT  = 0
	ROOM_TURN_CHOOSE_PASS = 1
)

type (
	RoomBet struct {
		PlayerID string `json:"player_id"`
		RoomID   string `json:"room_id"`
		Bet      int    `json:"bet"`
	}

	RoomTurn struct {
		PlayerID string `json:"player_id"`
		RoomID   string `json:"room_id"`
		Choosed  int    `json:"choosed"`
	}
)
