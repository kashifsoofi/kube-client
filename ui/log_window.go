package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kashifsoofi/kube-client/k8s"
)

func NewLogWindow(a fyne.App, client *k8s.Client, pod string) fyne.Window {
	w := a.NewWindow(pod + " Logs")

	logEntry := widget.NewMultiLineEntry()
	logEntry.Disable()

	logContainer := container.NewVScroll(
		logEntry,
	)
	logContainer.SetMinSize(fyne.NewSize(800, 600))

	content := container.NewVBox(
		container.NewHBox(
			widget.NewButton("Start", func() {

			}),
			widget.NewButton("Stop", func() {

			}),
		),
		logContainer,
	)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}
