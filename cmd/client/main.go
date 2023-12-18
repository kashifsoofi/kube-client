package main

import (
	"fmt"

	"github.com/kashifsoofi/kube-client/ui"
)

func main() {
	app, err := ui.NewApp()
	if err != nil {
		fmt.Printf("Panic: %+v\n", err)
	}
	app.Run()
}
