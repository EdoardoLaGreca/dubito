package main

import (
	"fyne.io/fyne/v2/app"
)

var username string = "Player"

func main() {
	a := app.New()
	w := a.NewWindow("Dubito")

	menuContainer := getMenuContainer(w)

	w.SetContent(menuContainer)
	w.ShowAndRun()
}
