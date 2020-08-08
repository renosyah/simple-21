package model

type (
	GameConfig struct {
		MaxPlayer    int `json:"max_player"`
		MaxRoom      int `json:"max_room"`
		StarterMoney int `json:"starter_money"`
	}
)
