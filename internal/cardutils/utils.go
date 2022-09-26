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

func RankByName(name string) (Rank, error) {
	var r Rank

	switch name {
	case "ace":
		r = Ace
	case "two":
		r = Two
	case "three":
		r = Three
	case "four":
		r = Four
	case "five":
		r = Five
	case "six":
		r = Six
	case "seven":
		r = Seven
	case "eight":
		r = Eight
	case "nine":
		r = Nine
	case "ten":
		r = Ten
	case "jack":
		r = Jack
	case "queen":
		r = Queen
	case "king":
		r = King
	default:
		return 0, fmt.Errorf("unknown rank in " + name)
	}

	return r, nil
}

// e.g. "five clubs" or "queen spades"
func CardByName(name string) (Card, error) {
	nameSp := strings.Fields(name)

	rankStr := nameSp[0]
	suitStr := nameSp[1]

	var c Card

	r, err := RankByName(rankStr)
	if err != nil {
		return Card{}, err
	}
	c.Rank = r

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

func RankToString(r Rank) string {
	var s string

	switch r {
	case Ace:
		s = "ace"
	case Two:
		s = "two"
	case Three:
		s = "three"
	case Four:
		s = "four"
	case Five:
		s = "five"
	case Six:
		s = "six"
	case Seven:
		s = "seven"
	case Eight:
		s = "eight"
	case Nine:
		s = "nine"
	case Ten:
		s = "ten"
	case Jack:
		s = "jack"
	case Queen:
		s = "queen"
	case King:
		s = "king"
	}

	return s
}

func CardToString(c Card) (name string) {
	name += RankToString(c.Rank) + " "

	switch c.Suit {
	case Clubs:
		name += "clubs"
	case Diamonds:
		name += "diamonds"
	case Hearts:
		name += "hearts"
	case Spades:
		name += "spades"
	}

	return
}

// CardsToString converts multiple cards to strings
func CardsToString(cards []Card) []string {
	cardsStr := make([]string, len(cards))

	for i, c := range cards {
		cardsStr[i] = CardToString(c)
	}

	return cardsStr
}
