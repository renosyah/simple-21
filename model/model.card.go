package model

import (
	"fmt"
)

var CARD_GROUPS = map[string]string{
	"Spade":   "./img/Spade.png",
	"Diamond": "./img/Diamond.png",
	"Heart":   "./img/Heart.png",
	"Club":    "./img/Club.png",
}

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

func NewCards(cgroups []string) []Card {
	var cards []Card

	groups := cgroups
	if len(cgroups) == 0 {
		for k, _ := range CARD_GROUPS {
			groups = append(groups, k)
		}
	}

	for _, g := range groups {
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
				Image: CARD_GROUPS[g],
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
