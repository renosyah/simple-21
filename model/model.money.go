package model

type (
	Money struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	PuchaseMoney struct {
		ID       string `json:"id"`
		PlayerID string `json:"player_id"`
	}
)

func ListMoney() []Money {
	return []Money{
		{
			ID:     "money-1",
			Name:   "Small Package",
			Amount: 100,
		},
		{
			ID:     "money-2",
			Name:   "Starter Package",
			Amount: 500,
		},
		{
			ID:     "money-3",
			Name:   "Medium Package",
			Amount: 750,
		},
		{
			ID:     "money-4",
			Name:   "Rich Gambler Package",
			Amount: 1000,
		},
	}
}
