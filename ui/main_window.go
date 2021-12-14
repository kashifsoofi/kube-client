package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kashifsoofi/kube-client/k8s"
)

var content *fyne.Container

func NewMainWindow(a fyne.App, client *k8s.Client) fyne.Window {
	w := a.NewWindow("Kubernetes Client")

	widgetNamespace := &widget.Select{}

	contexts, current := getContexts(client)
	widgetContext := widget.NewSelect(contexts, func(name string) {
		widgetNamespace.Options = switchContext(client, name)
		widgetNamespace.ClearSelected()
	})
	widgetContext.SetSelected(current)

	widgetResource := widget.NewSelect(getResources(), func(name string) {
		loadResources(name, client, widgetNamespace.Selected)
	})

	content = container.NewMax()

	mainContent := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Context"),
			widgetContext,
			widget.NewLabel("Namespace"),
			widgetNamespace,
			widget.NewLabel("Resource"),
			widgetResource,
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

func getContexts(client *k8s.Client) ([]string, string) {
	contexts := client.GetContexts()
	current := client.GetCurrentContext()

	return contexts, current
}

func switchContext(client *k8s.Client, ctx string) []string {
	client.SwitchContext(ctx)
	return getNamespaces(client)
}

func getNamespaces(client *k8s.Client) []string {
	namespaces, err := client.GetNamespaces()
	if err != nil {
		return []string{}
	}

	return namespaces
}

func getResources() []string {
	return []string{
		"Services",
		"Pods",
		"Deployments",
	}
}

func loadResources(name string, client *k8s.Client, ns string) {
	switch name {
	case "Pods":
		loadPods(client, ns)
	case "Services":
		loadServices(client, ns)
	case "Deployments":
		loadDeployments(client, ns)
	}
}

func loadPods(client *k8s.Client, ns string) {
	podNames, err := client.GetPods(ns)
	if err != nil {
		return
	}

	podCards := []fyne.CanvasObject{}
	for _, p := range podNames {
		podCard := widget.NewCard(
			p,
			"",
			container.NewHBox(
				widget.NewButton("Logs", func() {

				}),
				widget.NewButton("Scale", func() {

				}),
				widget.NewButton("Delete", func() {

				}),
			))
		podCards = append(podCards, podCard)
	}

	pods := container.NewVScroll(
		container.NewVBox(podCards...),
	)
	pods.SetMinSize(fyne.NewSize(640, 460))

	content.Objects = []fyne.CanvasObject{pods}
	content.Refresh()
}

func loadServices(client *k8s.Client, ns string) {
	serviceNames, err := client.GetServices(ns)
	if err != nil {
		return
	}

	serviceCards := []fyne.CanvasObject{}
	for _, s := range serviceNames {
		serviceCard := widget.NewCard(
			s,
			"",
			container.NewHBox(
				widget.NewButton("Delete", func() {

				}),
			))
		serviceCards = append(serviceCards, serviceCard)
	}

	services := container.NewVScroll(
		container.NewVBox(serviceCards...),
	)
	services.SetMinSize(fyne.NewSize(640, 460))

	content.Objects = []fyne.CanvasObject{services}
	content.Refresh()
}

func loadDeployments(client *k8s.Client, ns string) {
	deploymentNames, err := client.GetDeployments(ns)
	if err != nil {
		return
	}

	deploymentCards := []fyne.CanvasObject{}
	for _, dn := range deploymentNames {
		deploymentCard := widget.NewCard(
			dn,
			"",
			container.NewHBox(
				widget.NewButton("Port Forward", func() {

				}),
				widget.NewButton("Delete", func() {

				}),
			))
		deploymentCards = append(deploymentCards, deploymentCard)
	}

	deployments := container.NewVScroll(
		container.NewVBox(deploymentCards...),
	)
	deployments.SetMinSize(fyne.NewSize(640, 460))

	content.Objects = []fyne.CanvasObject{deployments}
	content.Refresh()
}
