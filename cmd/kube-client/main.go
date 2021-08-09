package main

import (
	"github.com/kashifsoofi/kube-client/ui"
)

func main() {
	app := ui.NewApp(store)
	app.Run()
}