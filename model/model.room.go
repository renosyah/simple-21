package model

type (
	Room struct {
		ID           string       `json:"id"`
		PlayerTurnID string       `json:"player_turn_id"`
		RoomPlayers  []RoomPlayer `json:"room_players"`
		Cards        []Card       `json:"-"`
	}
)
