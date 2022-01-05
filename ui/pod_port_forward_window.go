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

type PodPortForwardWindow struct {
	w    fyne.Window
	logs *widget.Label

	client *k8s.Client
	ns     string
	pn     string

	wg sync.WaitGroup
	// stopCh control the port forwarding lifecycle. When it gets closed the
	// port forward will terminate
	stopCh chan struct{}
}

func NewPodPortForwardWindow(a fyne.App, client *k8s.Client, ns, pod string) fyne.Window {
	w := a.NewWindow(pod + " Port Forward")

	logsLabel := widget.NewLabel("")

	scrollableContent := container.NewScroll(
		logsLabel,
	)
	scrollableContent.SetMinSize(fyne.NewSize(800, 600))
	scrollableContent.ScrollToBottom()

	localPortBinding := binding.NewString()
	localPortBinding.Set("8081")
	localPortEntry := widget.NewEntryWithData(localPortBinding)

	podPortBinding := binding.NewString()
	podPortBinding.Set("80")
	podPortEntry := widget.NewEntryWithData(podPortBinding)

	pfw := PodPortForwardWindow{
		w:      w,
		logs:   logsLabel,
		client: client,
		ns:     ns,
		pn:     pod,
		wg:     sync.WaitGroup{},
	}

	in := make(chan string)
	go pfw.refreshLogs(in)

	content := container.NewVBox(
		container.NewHBox(
			localPortEntry,
			podPortEntry,
			widget.NewButton("Start", func() {
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

				pfw.wg.Add(1)
				// readyCh communicate when the port forward is ready to get traffic
				readyCh := make(chan struct{})
				pfw.stopCh = make(chan struct{}, 1)
				localPort, _ := localPortBinding.Get()
				podPort, _ := podPortBinding.Get()
				go pfw.startPodPortForward(
					localPort,
					podPort,
					outWriter,
					errOutWriter,
					readyCh)

				// read ready channel
				<-readyCh

				go pfw.routeOutput(outReader, errOutReader, in)

				// pfw.wg.Wait()
			}),
			widget.NewButton("Stop", func() {
				if pfw.stopCh == nil {
					fmt.Println("already stopped")
					return
				}
				close(pfw.stopCh)
				pfw.wg.Done()
				logsLabel.Text += "Port forwarding is stopped for " + pod + ".\n"
			}),
		),
		scrollableContent,
	)

	w.SetContent(content)
	w.CenterOnScreen()
	return w
}

func (pfw PodPortForwardWindow) startPodPortForward(localPort, podPort string, outWriter, errOutWriter *os.File, readyCh chan struct{}) {

	err := pfw.client.PortForwardAPod(k8s.PortForwardAPodRequest{
		Namespace: pfw.ns,
		PodName:   pfw.pn,
		LocalPort: localPort,
		PodPort:   podPort,
		Out:       outWriter,
		ErrOut:    errOutWriter,
		StopCh:    pfw.stopCh,
		ReadyCh:   readyCh,
	})
	if err != nil {
		dialog.NewError(err, pfw.w)
	}
}

func (pfw PodPortForwardWindow) routeOutput(outReader, errOutReader *os.File, out chan<- string) {
	for {
		select {
		case <-pfw.stopCh:
			fmt.Println("Stopped")
			return

		default:
			output, err := readOutput(outReader)
			if err == io.EOF {
				fmt.Println("Output EOF")
				break
			}
			if err != nil {
				fmt.Printf("Output Error: %v\n", err)
				return
			}

			out <- output

			errOutput, err := readOutput(errOutReader)
			if err == io.EOF {
				fmt.Println("Error Output EOF")
				break
			}
			if err != nil {
				fmt.Printf("ErrorOutput Error: %v\n", err)
				return
			}

			out <- errOutput
		}
	}
}

func (pfw PodPortForwardWindow) refreshLogs(in <-chan string) {
	for {
		text := <-in
		pfw.logs.Text += text
		pfw.logs.Refresh()
	}
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
