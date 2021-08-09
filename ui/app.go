package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	app    fyne.App
	window fyne.Window
}

func NewApp() *App {
	app := App{
		app:   app.NewWithID("Kubernetes Client"),
	}
	app.window = NewMainWindow(app.app)

	return &app
}

func (a *App) Run() {
	a.window.ShowAndRun()
}