package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func changePages(content *fyne.Container, newContent *fyne.Container) {
	content.Objects = []fyne.CanvasObject{newContent}
	content.Refresh()
}

func makePlayContent() *fyne.Container {
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Fortnite Path (D:\\ProjectFortress\\)")

	return container.NewVBox(
		pathEntry,
	)
}

func makeExitContent(w fyne.Window) *fyne.Container {
	return container.NewVBox(
		widget.NewLabelWithStyle("Are you sure you want to exit?", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewButton("Yes", func() {
			w.Close()
		}),
	)
}
