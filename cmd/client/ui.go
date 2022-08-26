package main

import (
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/dialog"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func getSettingsContainer(w fyne.Window) *fyne.Container {
	lblUsername := widget.NewLabel("Username")
	lblAddress := widget.NewLabel("Server address")
	lblPort := widget.NewLabel("Server port")

	entUsername := widget.NewEntry()
	entAddress := widget.NewEntry()
	entPort := widget.NewEntry()

	return container.New(layout.NewGridLayout(2), lblUsername, entUsername, lblAddress, entAddress, lblPort, entPort)
}

func getWaitingRoomContainer(w fyne.Window, maxPlayers uint) *fyne.Container {
	lblJoined := widget.NewLabel("0/" + strconv.Itoa(int(maxPlayers)) + " joined")

	return container.New(layout.NewCenterLayout(), lblJoined)
}

func getGameContainer(w fyne.Window, players []string) *fyne.Container {
	cnvPlayers := make([]*canvas.Text, len(players))
	for i := range players {
		cnvPlayers[i] = canvas.NewText(players[i], color.RGBA{R: 200, G: 200, B: 200, A: 255})
	}

	return container.New(layout.NewVBoxLayout())
}

func getMenuContainer(w fyne.Window) *fyne.Container {
	btnNewGame := widget.NewButton("New game", func() {
		err := initConn()
		if err != nil {
			dialog.ShowError(err, w)
		}

		maxPlayers, err := requestMaxPlayers(conn)
		if err != nil {
			dialog.ShowError(err, w)
		}

		wrCont := getWaitingRoomContainer(w, maxPlayers)
		w.SetContent(wrCont)

		lblJoined := wrCont.Objects[0]

		var players []string

		// wait until all players joined
		for len(players) < int(maxPlayers) {
			requestPlayers(conn)
			time.Sleep(time.Duration(500 * time.Millisecond))
		}
	})

	btnResumeGame := widget.NewButton("Resume game", func() {

	})

	btnSettings := widget.NewButton("Settings", func() {

	})

	return container.New(layout.NewVBoxLayout(), btnNewGame, btnResumeGame, btnSettings)
}
