package cardutils

import (
	"fmt"
	"strings"
)

type Card struct {
	Suit Suit
	Rank Rank
}

// suits
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

// ranks
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// e.g. "five clubs" or "queen spades"
func CardByName(name string) (Card, error) {
	nameSp := strings.Fields(name)

	rankStr := nameSp[0]
	suitStr := nameSp[1]

	var c Card

	switch rankStr {
	case "ace":
		c.Rank = Ace
	case "two":
		c.Rank = Two
	case "three":
		c.Rank = Three
	case "four":
		c.Rank = Four
	case "five":
		c.Rank = Five
	case "six":
		c.Rank = Six
	case "seven":
		c.Rank = Seven
	case "eight":
		c.Rank = Eight
	case "nine":
		c.Rank = Nine
	case "ten":
		c.Rank = Ten
	case "jack":
		c.Rank = Jack
	case "queen":
		c.Rank = Queen
	case "king":
		c.Rank = King
	default:
		return Card{}, fmt.Errorf("unknown rank in " + name)
	}

	switch suitStr {
	case "clubs":
		c.Suit = Clubs
	case "diamonds":
		c.Suit = Diamonds
	case "hearts":
		c.Suit = Hearts
	case "spades":
		c.Suit = Spades
	default:
		return Card{}, fmt.Errorf("unknown suit in " + name)
	}

	return c, nil
}
