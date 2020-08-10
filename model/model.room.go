package model

const (
	ROOM_STATUS_USE     = 0
	ROOM_STATUS_NOT_USE = 1
)

type (
	Room struct {
		ID        string       `json:"id"`
		OwnerID   string       `json:"-"`
		Name      string       `json:"name"`
		Dealer    RoomPlayer   `json:"dealer"`
		Players   []RoomPlayer `json:"players"`
		Removable bool         `json:"removable"`
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
