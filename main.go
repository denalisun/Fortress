package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

func main() {
	fmt.Println("Hello, world!")

	a := app.New()
	w := a.NewWindow("Fortress")
	w.Resize(fyne.NewSize(300, 400))

	w.ShowAndRun()
}
