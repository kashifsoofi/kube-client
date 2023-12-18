package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewAddContextWindow(a fyne.App) fyne.Window {
	addContextContainer := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Context Name"),
			widget.NewEntry(),
		),
		container.NewHBox(
			widget.NewLabel("Namespace"),
			widget.NewEntry(),
		),
		container.NewHBox(
			widget.NewLabel("Cluster"),
			widget.NewEntry(),
		),
		container.NewHBox(
			widget.NewLabel("User"),
			widget.NewEntry(),
		),
	)

	w := a.NewWindow("Add Context")
	w.SetContent(addContextContainer)
	return w
}
