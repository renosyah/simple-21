package model

import (
	"fmt"
)

var CARD_GROUPS = []string{"Spade", "Diamond", "Heart", "Club"}
var CARD_GROUPS_URL = []string{"./img/Spade.png", "./img/Diamond.png", "./img/Heart.png", "./img/Club.png"}

type (
	Card struct {
		ID    string `json:"id"`
		Label string `json:"label"`
		Value int    `json:"value"`
		Group string `json:"group"`
		Show  bool   `json:"show"`
		Image string `json:"image"`
	}
)

func (c *Card) Copy(show bool) Card {
	return Card{
		ID:    c.ID,
		Label: c.Label,
		Value: c.Value,
		Group: c.Group,
		Show:  show,
		Image: c.Image,
	}
}
func (c *Card) CopyPointer() *Card {
	return &Card{
		ID:    c.ID,
		Label: c.Label,
		Value: c.Value,
		Group: c.Group,
		Show:  c.Show,
		Image: c.Image,
	}
}

func NewCards() []Card {
	var cards []Card
	for pos, g := range CARD_GROUPS {
		for i := 1; i <= 13; i++ {
			value := i

			// J,Q & K = 10
			if i == 11 || i == 12 || i == 13 {
				value = 10
			}

			cards = append(cards, Card{
				ID:    fmt.Sprintf("card-%s-%d", g, i),
				Label: GetLabel(fmt.Sprint(i)),
				Value: value,
				Group: g,
				Image: CARD_GROUPS_URL[pos],
				Show:  false,
			})
		}
	}
	return cards
}

func GetLabel(value string) string {
	switch value {
	case "1":
		return "A"
	case "11":
		return "J"
	case "12":
		return "Q"
	case "13":
		return "K"
	default:
		break
	}
	return value
}
