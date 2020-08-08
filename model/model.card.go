package model

import (
	"fmt"
)

var CARD_GROUPS = []string{"Spade", "Diamond", "Heart", "Club"}

type (
	Card struct {
		ID    string `json:"id"`
		Label string `json:"label"`
		Value int    `json:"value"`
		Group string `json:"group"`
		Show  bool   `json:"show"`
	}
)

func NewCards() []Card {
	var cards []Card
	for _, g := range CARD_GROUPS {
		for i := 1; i <= 12; i++ {
			cards = append(cards, Card{
				ID:    fmt.Sprintf("card-%s-%d", g, i),
				Label: GetLabel(fmt.Sprint(i)),
				Value: i,
				Group: g,
			})
		}
	}
	return cards
}

func GetLabel(value string) string {
	switch value {
	case "1":
		return "AS"
	case "10":
		return "J"
	case "11":
		return "Q"
	case "12":
		return "K"
	default:
		break
	}
	return value
}
