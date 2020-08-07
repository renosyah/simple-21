package model

import (
	"fmt"
)

type (
	Card struct {
		ID    string `json:"id"`
		Label string `json:"label"`
		Value string `json:"value"`
	}
)

func NewCards() []Card {
	var cards []Card
	for i := 1; i <= 12; i++ {
		cards = append(cards, Card{
			ID:    fmt.Sprint("CARD-", i),
			Label: GetLabel(fmt.Sprint(i)),
			Value: fmt.Sprint(i),
		})
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
