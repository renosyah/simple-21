package model

type (
	Room struct {
		ID           string       `json:"id"`
		PlayerTurnID string       `json:"player_turn_id"`
		Players      []PlayerRoom `json:"players"`
		Cards        []Card       `json:"-"`
	}
)
