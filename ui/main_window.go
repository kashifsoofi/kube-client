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
		loadResources(a, name, client, widgetNamespace.Selected)
	})

	subContent := container.NewVScroll(
		container.NewVBox(),
	)
	subContent.SetMinSize(fyne.NewSize(640, 480))
	content = container.NewCenter(subContent)

	mainContent := container.NewBorder(
		container.NewHBox(
			widget.NewLabel("Context"),
			widgetContext,
			widget.NewLabel("Namespace"),
			widgetNamespace,
			widget.NewLabel("Resource"),
			widgetResource,
		),
		nil,
		nil,
		nil,
		content,
	)

	w.SetContent(mainContent)
	w.CenterOnScreen()
	return w
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

func loadResources(a fyne.App, name string, client *k8s.Client, ns string) {
	var cards []fyne.CanvasObject
	switch name {
	case "Pods":
		cards = loadPods(a, client, ns)
	case "Services":
		cards = loadServices(client, ns)
	case "Deployments":
		cards = loadDeployments(client, ns)
	}

	cardsContainer := container.NewVScroll(
		container.NewGridWrap(fyne.NewSize(450, 100), cards...),
	)

	cardsContainer.SetMinSize(fyne.NewSize(920, 600))

	content.Objects = []fyne.CanvasObject{cardsContainer}
	content.Refresh()
}

func loadPods(a fyne.App, client *k8s.Client, ns string) []fyne.CanvasObject {
	podNames, err := client.GetPods(ns)
	if err != nil {
		return []fyne.CanvasObject{}
	}

	podCards := []fyne.CanvasObject{}
	for _, p := range podNames {
		podCard := widget.NewCard(
			"",
			p,
			container.NewHBox(
				widget.NewButton("Logs", func() {
					lw := NewLogWindow(a, client, p)
					lw.Show()
				}),
				widget.NewButton("Delete", func() {

				}),
			))
		podCards = append(podCards, podCard)
	}

	return podCards
}

func loadServices(client *k8s.Client, ns string) []fyne.CanvasObject {
	serviceNames, err := client.GetServices(ns)
	if err != nil {
		return []fyne.CanvasObject{}
	}

	serviceCards := []fyne.CanvasObject{}
	for _, s := range serviceNames {
		serviceCard := widget.NewCard(
			"",
			s,
			container.NewHBox(
				widget.NewButton("Delete", func() {

				}),
			))
		serviceCards = append(serviceCards, serviceCard)
	}

	return serviceCards
}

func loadDeployments(client *k8s.Client, ns string) []fyne.CanvasObject {
	deploymentNames, err := client.GetDeployments(ns)
	if err != nil {
		return []fyne.CanvasObject{}
	}

	deploymentCards := []fyne.CanvasObject{}
	for _, dn := range deploymentNames {
		deploymentCard := widget.NewCard(
			"",
			dn,
			container.NewHBox(
				widget.NewButton("Port Forward", func() {

				}),
				widget.NewButton("Scale", func() {

				}),
				widget.NewButton("Delete", func() {

				}),
			))
		deploymentCards = append(deploymentCards, deploymentCard)
	}

	return deploymentCards
}
