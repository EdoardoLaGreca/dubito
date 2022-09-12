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
var imageAssets embed.FS

func getCardFilename(c cardutils.Card) (filename string) {
	filename = "card_"

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

	return
}

func GetCardAsset(c cardutils.Card) (image.Image, error) {
	filename := getCardFilename(c)

	content, err := imageAssets.ReadFile("cards/" + filename)
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
	content, err := imageAssets.ReadFile("decks/deck_" + strconv.Itoa(style) + ".png")
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}

	return img, nil
}
