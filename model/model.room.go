package model

type (
	Room struct {
		ID      string   `json:"id"`
		OwnerID string   `json:"-"`
		Name    string   `json:"name"`
		Players []Player `json:"players"`
		Round   int      `json:"round"`
	}
	AddRoom struct {
		HostID  string   `json:"host_id"`
		Name    string   `json:"name"`
		Players []Player `json:"players"`
	}
	DeleteRoom struct {
		ID       string `json:"id"`
		PlayerID string `json:"player_id"`
	}
)
