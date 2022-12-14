package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/EdoardoLaGreca/dubito/assets"
	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
)

var deckStyle int = 1
var selectedCards []cardutils.Card = make([]cardutils.Card, 0)

// true if selectedCards contains card
func selectedCardsContains(card cardutils.Card) bool {
	for _, c := range selectedCards {
		if c == card {
			return true
		}
	}

	return false
}

// remove card from selectedCards if it exists, otherwise do nothing
func removeSelectedCard(card cardutils.Card) {
	for i, c := range selectedCards {
		if c == card {
			selectedCards = append(selectedCards[:i], selectedCards[i+1:]...)
			break
		}
	}
}

func getSettingsContainer(w fyne.Window) *fyne.Container {
	lblUsername := widget.NewLabel("Username")
	entUsername := widget.NewEntry()
	entUsername.OnChanged = func(value string) {
		username = value
	}
	entUsername.Text = username

	lblAddress := widget.NewLabel("Server address")
	entAddress := widget.NewEntry()
	entAddress.OnChanged = func(value string) {
		serverAddress = value
	}
	entAddress.Text = serverAddress

	lblPort := widget.NewLabel("Server port")
	entPort := widget.NewEntry()
	entPort.OnChanged = func(value string) {
		port, err := strconv.Atoi(value)
		if err != nil || port >= 1<<16 {
			dialog.ShowError(fmt.Errorf("invalid port number"), w)
		} else {
			serverPort = uint16(port)
		}
	}
	entPort.Text = strconv.Itoa(int(serverPort))

	lblDeckStyle := widget.NewLabel("Deck style")
	cmbDeckStyle := widget.NewSelect(make([]string, 0), func(value string) {
		styleNumber := strings.Fields(value)[1]
		deckStyle, _ = strconv.Atoi(styleNumber)
	})

	for i := 1; i <= 6; i++ {
		cmbDeckStyle.Options = append(cmbDeckStyle.Options, "Style "+strconv.Itoa(i))
	}

	cmbDeckStyle.SetSelected("Style " + strconv.Itoa(int(deckStyle)))

	btnBack := widget.NewButton("Back", func() {
		w.SetContent(getMenuContainer(w))
	})

	return container.New(layout.NewGridLayout(2), lblUsername, entUsername, lblAddress, entAddress, lblPort, entPort, lblDeckStyle, cmbDeckStyle, btnBack)
}

func getWaitingRoomContainer(w fyne.Window, maxPlayers uint) *fyne.Container {
	lblJoined := widget.NewLabel("0/" + strconv.Itoa(int(maxPlayers)) + " joined")

	return container.New(layout.NewCenterLayout(), lblJoined)
}

func updateJoinedCount(label *widget.Label, joinedPlayers, maxPlayers uint) {
	label.SetText(strconv.Itoa(int(joinedPlayers)) + "/" + strconv.Itoa(int(maxPlayers)) + " joined")
}

// create a cards container
func newCardsCont(w fyne.Window, lblSelectedCards *widget.Label, cards []cardutils.Card) *fyne.Container {
	// use a grid layout of 1 row to divide the horizontal space in equal parts
	cardsCont := container.New(layout.NewGridWrapLayout(fyne.NewSize(40.0, 80.0)))

	for _, c := range cards {
		img, err := assets.GetCardAsset(c)
		if err != nil {
			dialog.ShowError(err, w)
			break
		}

		canvasImage := canvas.NewImageFromImage(img)
		canvasImage.FillMode = canvas.ImageFillContain

		rectSelected := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 255})

		// save the card in a new variable so that the function literal doesn't reference the variable updated by the loop. otherwise it would
		// cause function literals to all have the same value, which would always be the last card of the loop
		currentCard := c

		imgButton := widget.NewButton("", func() {
			if selectedCardsContains(currentCard) {
				removeSelectedCard(currentCard)
				rectSelected.FillColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
				rectSelected.Refresh()
			} else {
				if len(selectedCards) < 4 {
					selectedCards = append(selectedCards, currentCard)
					rectSelected.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
					rectSelected.Refresh()
				} else {
					dialog.ShowError(fmt.Errorf("you selected 4 cards already"), w)
				}
			}
			lblSelectedCards.Text = fmt.Sprintf("You selected %d cards: %s", len(selectedCards), strings.Join(cardutils.CardsToString(selectedCards), ", "))
			lblSelectedCards.Refresh()
		})

		// button with a rectangle and an image on top of it
		clickableImg := container.NewMax(imgButton, rectSelected, canvasImage)

		cardsCont.Add(clickableImg)
	}

	return cardsCont
}

func getGameContainer(w fyne.Window, players []string, cards []cardutils.Card) *fyne.Container {
	cnvPlayers := make([]fyne.CanvasObject, len(players))
	for i := range players {
		cnvPlayers[i] = canvas.NewText(players[i], color.RGBA{R: 200, G: 200, B: 200, A: 255})
	}

	// Set first player
	cnvPlayers[0].(*canvas.Text).Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}

	playersCont := container.New(layout.NewHBoxLayout(), cnvPlayers...)

	// initially, the last card is the deck style
	img, err := assets.GetDeckAsset(deckStyle)
	if err != nil {
		dialog.ShowError(err, w)
	}

	lastCardPlaced := canvas.NewImageFromImage(img)
	lastCardPlaced.FillMode = canvas.ImageFillContain
	lastCardPlaced.SetMinSize(fyne.NewSize(100.0, 200.0))
	lastCardCont := container.New(layout.NewBorderLayout(nil, nil, nil, nil), lastCardPlaced)

	lblSelectedCards := widget.NewLabel("You selected 0 cards")

	cardsCont := newCardsCont(w, lblSelectedCards, cards)

	btnPlace := widget.NewButton("Place cards", func() {
		if len(selectedCards) == 0 {
			return
		}

		ok, err := requestPlaceCards(selectedCards)
		if err != nil {
			dialog.ShowError(err, w)
		} else if !ok {
			dialog.ShowError(fmt.Errorf("unable to place cards"), w)
		}
	})

	btnDubito := widget.NewButton("Dubito!", func() {
		tableCards, err := requestDubito()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		if tableCards == nil {
			dialog.ShowInformation("Correct!", "You doubted well! :)", w)
		} else {
			dialog.ShowInformation("Wrong!", "You doubted wrong... :(", w)

			// add cards
			cardsCont = newCardsCont(w, lblSelectedCards, append(cards, tableCards...))
			cardsCont.Refresh()
		}

	})

	btnPlace.Hide()

	btnLeave := widget.NewButton("Leave", func() {
		requestLeave()
		w.SetContent(getMenuContainer(w))
	})

	return container.New(layout.NewVBoxLayout(), playersCont, lastCardCont, cardsCont, lblSelectedCards, btnPlace, btnDubito, btnLeave)
}

func getMenuContainer(w fyne.Window) *fyne.Container {
	btnNewGame := widget.NewButton("New game", func() {
		newGame(w)
	})

	btnSettings := widget.NewButton("Settings", func() {
		w.SetContent(getSettingsContainer(w))
	})

	return container.New(layout.NewGridLayoutWithColumns(1), btnNewGame, btnSettings)
}

// show a dialog containing the error (if not nil) and load the main menu container
func backToMainMenu(w fyne.Window, err error) {
	if err != nil {
		dialog.ShowError(err, w)
	}
	w.SetContent(getMenuContainer(w))
}

// connection closing handler, it does thing when the connection is lost
func connClosingHandler(w fyne.Window) {
	// the connection is lost
	<-closeChan

	dialog.ShowError(fmt.Errorf("connection lost"), w)
	w.SetContent(getMenuContainer(w))
}

func newGame(w fyne.Window) {
	err := initConn()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	go connClosingHandler(w)

	err = requestJoin()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	w.SetTitle("Dubito | in game as " + username)

	maxPlayers, err := requestMaxPlayers()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	wrCont := getWaitingRoomContainer(w, maxPlayers)
	w.SetContent(wrCont)

	lblJoined := wrCont.Objects[0].(*widget.Label)

	var players []string

	// wait until all players joined
	for len(players) < int(maxPlayers) {
		var err error
		players, err = requestPlayers()
		if err != nil {
			backToMainMenu(w, err)
			return
		}

		updateJoinedCount(lblJoined, uint(len(players)), maxPlayers)
		time.Sleep(time.Duration(200 * time.Millisecond))
	}

	cards, err := requestCards()
	if err != nil {
		backToMainMenu(w, err)
		return
	}

	gameCont := getGameContainer(w, players, cards)
	w.SetContent(gameCont)

	cnvLastCard := gameCont.Objects[1].(*fyne.Container).Objects[0].(*canvas.Image)
	btnPlace := gameCont.Objects[4].(*widget.Button)

	// goroutine to perform actions based on turns
	go func(w fyne.Window) {
		for {
			ud, err := requestUpdate()
			if err != nil {
				continue
			}

			if ud.gameOver {
				if ud.playerWon {
					dialog.ShowInformation("You won!", "Congrats, you won this game! :)", w)
				} else {
					dialog.ShowInformation("You lost...", "You lost this game. :(", w)
				}
			} else {
				if ud.playerTurn {
					btnPlace.Show()
				} else {
					btnPlace.Hide()
				}
			}

			// update last card
			newLastCard := cardutils.Card{Suit: cardutils.Spades, Rank: ud.cardRank}
			newLastCardAsset, err := assets.GetCardAsset(newLastCard)
			if err != nil {
				dialog.ShowError(err, w)
				continue
			}
			cnvLastCard.Image = newLastCardAsset
			cnvLastCard.Refresh()

			time.Sleep(200 * time.Millisecond)
		}
	}(w)
}
