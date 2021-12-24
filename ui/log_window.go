package ui

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kashifsoofi/kube-client/k8s"
)

type LogWindow struct {
	w         fyne.Window
	container *container.Scroll
	logs      *widget.Label

	client *k8s.Client
	ns     string
	pn     string
	stopCh chan struct{}
}

var logWindow LogWindow

func NewLogWindow(a fyne.App, client *k8s.Client, ns, pod string) fyne.Window {
	w := a.NewWindow(pod + " Logs")

	logEntry := widget.NewLabel("")

	logContainer := container.NewScroll(
		logEntry,
	)
	logContainer.SetMinSize(fyne.NewSize(800, 600))
	logContainer.ScrollToBottom()

	logWindow = LogWindow{
		w:         w,
		container: logContainer,
		logs:      logEntry,
		client:    client,
		ns:        ns,
		pn:        pod,
	}

	content := container.NewVBox(
		container.NewHBox(
			widget.NewButton("Start", func() {
				logWindow.stopCh = make(chan struct{}, 1)

				go logWindow.startLog()
			}),
			widget.NewButton("Stop", func() {
				close(logWindow.stopCh)
			}),
		),
		logContainer,
	)

	logWindow.w.SetContent(content)
	logWindow.w.CenterOnScreen()
	return w
}

func (lw LogWindow) startLog() {
	stream, err := lw.client.GetPodLogStream(lw.ns, lw.pn)
	if err != nil {
		dialog.NewError(err, lw.w)
		return
	}
	defer stream.Close()

	for {
		select {
		case <-lw.stopCh:
			return

		default:
			buf := make([]byte, 2000)
			numBytes, err := stream.Read(buf)
			if numBytes == 0 {
				continue
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				dialog.NewError(err, lw.w)
				return
			}

			message := string(buf[:numBytes])
			lw.logs.Text += message
			lw.logs.Refresh()
			lw.container.ScrollToBottom()
		}
	}
}
