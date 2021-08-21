package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Kubernetes Client")

	hello := widget.NewLabel("Hello Fyne!")
	content := container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}
