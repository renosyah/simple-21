package model

type (
	Room struct {
		ID           string        `json:"id"`
		OwnerID      string        `json:"-"`
		Name         string        `json:"name"`
		PlayerTurnID string        `json:"player_turn_id"`
		Dealer       *RoomPlayer   `json:"dealer"`
		RoomPlayers  []*RoomPlayer `json:"room_players"`
		Round        int           `json:"round"`
		Cards        []Card        `json:"-"`
	}

	DeleteRoom struct {
		ID       string `json:"id"`
		PlayerID string `json:"player_id"`
	}
)
