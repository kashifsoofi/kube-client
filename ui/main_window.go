package ui

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kashifsoofi/kube-client/k8s"
	"k8s.io/client-go/util/homedir"
)

var content fyne.CanvasObject

func NewMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Kubernetes Client")

	widgetNamespace := &widget.Select{}

	contexts, current := getContexts()
	widgetContext := widget.NewSelect(contexts, func(name string) {
		widgetNamespace.Options = getNamespaces()
	})
	widgetContext.SetSelected(current)

	content = container.NewVBox(
		widget.NewLabel("Main Window"),
	)

	mainContent := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Namespace"),
			widgetNamespace,
			widget.NewLabel("Context"),
			widgetContext,
		),
		container.NewHBox(
			makeNav(),
			container.NewCenter(content),
		),
	)

	w.SetContent(mainContent)
	w.CenterOnScreen()
	return w
}

func makeNav() fyne.CanvasObject {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return []string{"Cluster", "Nodes", "Workloads", "Configuration", "Network", "Storage", "Namespaces"}
		},
		IsBranch: func(uid string) bool {
			return false
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(uid)
		},
		OnSelected: func(uid string) {
			setContent(uid)
		},
	}

	return tree
}

func setContent(text string) {
	content = container.NewVBox(
		widget.NewLabel("Main Window"),
	)
	// content.Refresh()
}

func getContexts() ([]string, string) {
	home := homedir.HomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")
	contexts, current, err := k8s.GetContexts(kubeConfigPath)
	if err != nil {
		return nil, ""
	}

	return contexts, current
}

func getNamespaces(context string) []string {
	return []string{
		"namespace 1",
		"namespace 2",
		"namespace 3",
	}
}
