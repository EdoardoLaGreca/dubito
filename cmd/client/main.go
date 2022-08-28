package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2/app"
)

// suits
type suit int

const (
	clubs suit = iota
	diamonds
	hearts
	spades
)

// ranks
type rank int

const (
	ace rank = iota
	two
	three
	four
	five
	six
	seven
	eight
	nine
	ten
	jack
	queen
	king
)

type card struct {
	suit suit
	rank rank
}

var username string = "Player"
var serverAddress string = ""
var serverPort uint16 = 9876

// e.g. "five clubs" or "queen spades"
func cardByName(name string) (card, error) {
	nameSp := strings.Fields(name)

	rankStr := nameSp[0]
	suitStr := nameSp[1]

	var c card

	switch rankStr {
	case "ace":
		c.rank = ace
	case "two":
		c.rank = two
	case "three":
		c.rank = three
	case "four":
		c.rank = four
	case "five":
		c.rank = five
	case "six":
		c.rank = six
	case "seven":
		c.rank = seven
	case "eight":
		c.rank = eight
	case "nine":
		c.rank = nine
	case "ten":
		c.rank = ten
	case "jack":
		c.rank = jack
	case "queen":
		c.rank = queen
	case "king":
		c.rank = king
	default:
		return card{}, fmt.Errorf("unknown rank in " + name)
	}

	switch suitStr {
	case "clubs":
		c.suit = clubs
	case "diamonds":
		c.suit = diamonds
	case "hearts":
		c.suit = hearts
	case "spades":
		c.suit = spades
	default:
		return card{}, fmt.Errorf("unknown suit in " + name)
	}

	return c, nil
}

func main() {
	a := app.New()
	w := a.NewWindow("Dubito")

	menuContainer := getMenuContainer(w)

	w.SetContent(menuContainer)
	w.ShowAndRun()
}
