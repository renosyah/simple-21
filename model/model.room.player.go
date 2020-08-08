package model

const (
	AS_VALUE_ELEVEN = 0
	AS_VALUE_ONE    = 1

	PLAYER_STATUS_IN   = 0
	PLAYER_STATUS_OUT  = 1
	PLAYER_STATUS_BUST = 2
)

type (
	RoomPlayer struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Bet    int    `json:"bet"`
		Cards  []Card `json:"cards"`
		Total  int    `json:"total"`
		Status int    `json:"status"`
	}
)

func (p *RoomPlayer) SumUpTotal() {
	t := 0
	for _, i := range p.Cards {
		t += i.Value
	}
	p.Total = t
}

func (p *RoomPlayer) ChangeAsValue(flag int) {
	cards := []Card{}
	for _, i := range p.Cards {
		if i.Label == "A" && flag == AS_VALUE_ELEVEN {
			i.Value = 11
		}
		if i.Label == "A" && flag == AS_VALUE_ONE {
			i.Value = 1
		}
		cards = append(cards, i)
	}
	p.Cards = cards
	p.SumUpTotal()
}
