package assets

import (
	"bytes"
	"embed"
	"image"
	"image/png"
	"strconv"

	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
)

//go:embed cards decks
var assets embed.FS

func GetCardAsset(c cardutils.Card) (image.Image, error) {
	filename := "card_"

	switch c.Suit {
	case cardutils.Clubs:
		filename += "c"
	case cardutils.Diamonds:
		filename += "d"
	case cardutils.Hearts:
		filename += "h"
	case cardutils.Spades:
		filename += "s"
	}

	switch c.Rank {
	case cardutils.Ace:
		filename += "a"
	case cardutils.Jack:
		filename += "j"
	case cardutils.Queen:
		filename += "q"
	case cardutils.King:
		filename += "k"
	default:
		filename += strconv.Itoa(int(c.Rank))
	}

	filename += ".png"

	content, err := assets.ReadFile("cards/" + filename)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func GetDeckAsset(style int) (image.Image, error) {
	content, err := assets.ReadFile("decks/deck_" + strconv.Itoa(style) + ".png")
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}

	return img, nil
}
