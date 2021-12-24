package ui

import (
	"fmt"
	"io"
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
	w := a.NewWindow(pod + " Port Forward")

	logsLabel := widget.NewLabel("")

	scrollableContent := container.NewScroll(
		logsLabel,
	)
	scrollableContent.SetMinSize(fyne.NewSize(800, 600))
	scrollableContent.ScrollToBottom()

	var wg sync.WaitGroup

	// stopCh control the port forwarding lifecycle. When it gets closed the
	// port forward will terminate
	var stopCh chan struct{}

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
				wg.Add(1)
				stopCh = make(chan struct{}, 1)
				localPort, _ := localPortBinding.Get()
				podPort, _ := podPortBinding.Get()
				go startPodPortForward(
					client,
					ns,
					pod,
					localPort,
					podPort,
					w,
					wg,
					stopCh,
					logsLabel)
			}),
			widget.NewButton("Stop", func() {
				close(stopCh)
				wg.Done()
				logsLabel.Text += "Port forwarding is stopped for " + pod + ".\n"
			}),
		),
		scrollableContent,
	)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}

func startPodPortForward(client *k8s.Client, ns, pod, localPort, podPort string, w fyne.Window, wg sync.WaitGroup, stopCh chan struct{}, logs *widget.Label) {
	// readyCh communicate when the port forward is ready to get traffic
	readyCh := make(chan struct{})
	// stream is used to tell the port forwarder where to place its output or
	// where to expect input if needed. For the port forwarding we just need
	// the output eventually
	outReader, outWriter, err := os.Pipe()
	if err != nil {
		return
	}

	errOutReader, errOutWriter, err := os.Pipe()
	if err != nil {
		return
	}

	go func() {
		// PortForward the pod specified from its port 9090 to the local port
		// 8080
		err := client.PortForwardAPod(k8s.PortForwardAPodRequest{
			Namespace: ns,
			PodName:   pod,
			LocalPort: localPort,
			PodPort:   podPort,
			Out:       outWriter,
			ErrOut:    errOutWriter,
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

	go func() {
		for {
			select {
			case <-stopCh:
				fmt.Println("Stopped")
				return

			default:
				fmt.Println("Not Stopped yet")
				output, err := readOutput(outReader)
				if err == io.EOF {
					fmt.Println("Output EOF")
					break
				}
				if err != nil {
					fmt.Printf("Output Error: %v\n", err)
					return
				}

				logs.Text += output

				errOutput, err := readOutput(errOutReader)
				if err == io.EOF {
					fmt.Println("Error Output EOF")
					break
				}
				if err != nil {
					fmt.Printf("ErrorOutput Error: %v\n", err)
					return
				}

				logs.Text += errOutput
				logs.Refresh()
			}
		}
	}()

	wg.Wait()
}

func readOutput(reader *os.File) (string, error) {
	buf := make([]byte, 1000)
	numBytes, err := reader.Read(buf)
	if numBytes == 0 {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	output := string(buf[:numBytes])
	return output, nil
}
