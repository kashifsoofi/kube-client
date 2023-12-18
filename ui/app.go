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

func NewApp() (*App, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	app := App{
		app:    app.NewWithID("Kubernetes Client"),
		client: client,
	}
	app.window = NewMainWindow(app.app, app.client)

	return &app, nil
}

func (a *App) Run() {
	a.window.ShowAndRun()
}
