package main

import (
	"fyne.io/fyne/v2/app"
)

var username string = "Player"
var serverAddress string = "localhost"
var serverPort uint16 = 9876

func main() {
	a := app.New()
	w := a.NewWindow("Dubito")

	menuContainer := getMenuContainer(w)

	w.SetContent(menuContainer)
	w.ShowAndRun()
}
