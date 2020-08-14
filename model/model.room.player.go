package model

const (
	AS_VALUE_ELEVEN = 0
	AS_VALUE_ONE    = 1

	PLAYER_STATUS_SPECTATE    = -1
	PLAYER_STATUS_INVITED     = 0
	PLAYER_STATUS_SET_BET     = 1
	PLAYER_STATUS_IDLE        = 2
	PLAYER_STATUS_AT_TURN     = 3
	PLAYER_STATUS_FINISH_TURN = 4
	PLAYER_STATUS_OUT         = 5
	PLAYER_STATUS_BUST        = 6
	PLAYER_STATUS_REWARDED    = 7
	PLAYER_STATUS_LOSE        = 8
)

type (
	RoomPlayer struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Bet       int    `json:"bet"`
		Money     int    `json:"money"`
		Cards     []Card `json:"cards"`
		TotalShow int    `json:"total_show"`
		Total     int    `json:"total"`
		Status    int    `json:"status"`
		IsOnline  bool   `json:"is_online"`
		TurnOrder int    `json:"-"`
		IsBot     bool   `json:"-"`
	}
)

func (c *RoomPlayer) Copy() RoomPlayer {
	return RoomPlayer{
		ID:        c.ID,
		Name:      c.Name,
		Bet:       c.Bet,
		Money:     c.Money,
		Cards:     c.Cards,
		TotalShow: c.TotalShow,
		Total:     c.Total,
		Status:    c.Status,
		IsOnline:  c.IsOnline,
		TurnOrder: c.TurnOrder,
	}
}

func (p *RoomPlayer) SumUpTotal() {
	t := 0
	tshow := 0
	for _, i := range p.Cards {
		if i.Show {
			tshow += i.Value
		}
		t += i.Value
	}
	p.Total = t
	p.TotalShow = tshow
}

func (p *RoomPlayer) ShowAllCard() {
	nc := []Card{}
	for _, c := range p.Cards {
		nc = append(nc, c.Copy(true))
	}

	p.Cards = nc
	p.SumUpTotal()
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
