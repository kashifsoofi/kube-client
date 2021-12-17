package ui

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kashifsoofi/kube-client/k8s"
)

func NewLogWindow(a fyne.App, client *k8s.Client, ns, pod string) fyne.Window {
	w := a.NewWindow(pod + " Logs")

	logEntry := widget.NewMultiLineEntry()
	logEntry.Disable()

	logContainer := container.NewVScroll(
		logEntry,
	)
	logContainer.SetMinSize(fyne.NewSize(800, 600))

	var logStream io.ReadCloser
	content := container.NewVBox(
		container.NewHBox(
			widget.NewButton("Start", func() {
				stream, err := client.GetPodLogStream(ns, pod)
				if err != nil {
					dialog.NewError(err, w)
					return
				}
				defer stream.Close()

				logStream = stream

				for {
					buf := make([]byte, 2000)
					numBytes, err := stream.Read(buf)
					if numBytes == 0 {
						continue
					}
					if err == io.EOF {
						break
					}
					if err != nil {
						dialog.NewError(err, w)
						return
					}

					message := string(buf[:numBytes])
					logEntry.Text += message
					logEntry.Refresh()
				}
			}),
			widget.NewButton("Stop", func() {
				logStream.Close()
			}),
		),
		logContainer,
	)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}
