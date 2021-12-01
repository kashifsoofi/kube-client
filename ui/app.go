package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/kashifsoofi/kube-client/k8s"
)

type App struct {
	app    fyne.App
	window fyne.Window
	client *k8s.Client
}

func NewApp() *App {
	client, _ := k8s.NewClient()

	app := App{
		app:    app.NewWithID("Kubernetes Client"),
		client: client,
	}
	app.window = NewMainWindow(app.app, app.client)

	return &app
}

func (a *App) Run() {
	a.window.ShowAndRun()
}
