package ui

import (
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kashifsoofi/kube-client/k8s"
)

func NewPodPortForwardWindow(a fyne.App, client *k8s.Client, ns, pod string) fyne.Window {
	w := a.NewWindow(pod + " Logs")

	logEntry := widget.NewMultiLineEntry()
	logEntry.Disable()

	logContainer := container.NewVScroll(
		logEntry,
	)
	logContainer.SetMinSize(fyne.NewSize(800, 600))

	// stopCh control the port forwarding lifecycle. When it gets closed the
	// port forward will terminate
	stopCh := make(chan struct{}, 1)

	localPortBinding := binding.NewString()
	localPortBinding.Set("8081")
	localPortEntry := widget.NewEntryWithData(localPortBinding)

	podPortBinding := binding.NewString()
	podPortBinding.Set("80")
	podPortEntry := widget.NewEntryWithData(podPortBinding)

	content := container.NewVBox(
		container.NewHBox(
			localPortEntry,
			podPortEntry,
			widget.NewButton("Start", func() {
				localPort, _ := localPortBinding.Get()
				podPort, _ := podPortBinding.Get()
				go startPodPortForward(
					client,
					ns,
					pod,
					localPort,
					podPort,
					w,
					stopCh)
			}),
			widget.NewButton("Stop", func() {
				close(stopCh)
			}),
		),
		logContainer,
	)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}

func startPodPortForward(client *k8s.Client, ns, pod, localPort, podPort string, w fyne.Window, stopCh chan struct{}) {
	var wg sync.WaitGroup
	wg.Add(1)

	// readyCh communicate when the port forward is ready to get traffic
	readyCh := make(chan struct{})
	// stream is used to tell the port forwarder where to place its output or
	// where to expect input if needed. For the port forwarding we just need
	// the output eventually
	out := os.Stdout
	errOut := os.Stderr

	go func() {
		// PortForward the pod specified from its port 9090 to the local port
		// 8080
		err := client.PortForwardAPod(k8s.PortForwardAPodRequest{
			Namespace: ns,
			PodName:   pod,
			LocalPort: localPort,
			PodPort:   podPort,
			Out:       out,
			ErrOut:    errOut,
			StopCh:    stopCh,
			ReadyCh:   readyCh,
		})
		if err != nil {
			dialog.NewError(err, w)
		}
	}()

	select {
	case <-readyCh:
		break
	}
	println("Port forwarding is ready to get traffic. have fun!")

	wg.Wait()
}
