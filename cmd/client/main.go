package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var username string = "Player"

func main() {
	a := app.New()
	w := a.NewWindow("Dubito")
	w.Resize(fyne.NewSize(450, 200))

	menuContainer := getMenuContainer(w)

	w.SetContent(menuContainer)
	w.ShowAndRun()
}
