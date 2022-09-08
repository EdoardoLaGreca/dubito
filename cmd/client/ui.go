package main

import (
	"fmt"
	"image"
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

func getGameContainer(w fyne.Window, players []string, cards []cardutils.Card) *fyne.Container {
	cnvPlayers := make([]fyne.CanvasObject, len(players))
	for i := range players {
		cnvPlayers[i] = canvas.NewText(players[i], color.RGBA{R: 200, G: 200, B: 200, A: 255})
	}

	// Set first player
	cnvPlayers[0].(*canvas.Text).Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}

	playersCont := container.New(layout.NewHBoxLayout(), cnvPlayers...)

	lastCardPlaced := canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 390, 606)))

	myCards := make([]fyne.CanvasObject, len(cards))

	for i, c := range cards {
		img, err := assets.GetCardAsset(c)
		if err != nil {
			dialog.ShowError(err, w)
			break
		}

		myCards[i] = canvas.NewImageFromImage(img)
	}

	cardsCont := container.New(layout.NewHBoxLayout(), myCards...)

	return container.New(layout.NewVBoxLayout(), playersCont, lastCardPlaced, cardsCont)
}

func getMenuContainer(w fyne.Window) *fyne.Container {
	btnNewGame := widget.NewButton("New game", func() {
		newGame(w)
	})

	btnSettings := widget.NewButton("Settings", func() {
		w.SetContent(getSettingsContainer(w))
	})

	return container.New(layout.NewVBoxLayout(), btnNewGame, btnSettings)
}

// show a dialog containing the error (if not nil) and load the main menu container
func backToMainMenu(w fyne.Window, err error) {
	if err != nil {
		dialog.ShowError(err, w)
	}
	w.SetContent(getMenuContainer(w))
}

func newGame(w fyne.Window) {
	err := initConn()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	err = requestJoin(conn)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	maxPlayers, err := requestMaxPlayers(conn)
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
		players, err = requestPlayers(conn)
		if err != nil {
			backToMainMenu(w, err)
			return
		}

		updateJoinedCount(lblJoined, uint(len(players)), maxPlayers)
		time.Sleep(time.Duration(200 * time.Millisecond))
	}

	cards, err := requestCards(conn)
	if err != nil {
		backToMainMenu(w, err)
		return
	}

	gameCont := getGameContainer(w, players, cards)
	w.SetContent(gameCont)
}
