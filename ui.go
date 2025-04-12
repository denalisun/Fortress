package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func changePages(content *fyne.Container, newContent *fyne.Container) {
	content.Objects = []fyne.CanvasObject{newContent}
	content.Refresh()
}

func makePlayContent(settings *LauncherSettings) *fyne.Container {
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Fortnite Path (D:\\ProjectFortress\\)")
	pathEntry.Text = settings.FortniteInstallPath
	pathEntry.OnChanged = func(text string) {
		fmt.Println(text)
		settings.FortniteInstallPath = text
	}

	launchBtn := widget.NewButton("Launch Fortress", func() {
		go launchGame(settings)
	})

	return container.NewVBox(
		pathEntry,
		launchBtn,
	)
}

func makeOptionsContent(settings *LauncherSettings) *fyne.Container {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("E-mail/Username")
	usernameEntry.Text = settings.Username
	usernameEntry.OnChanged = func(text string) {
		settings.Username = text
	}

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")
	passwordEntry.Text = settings.Password
	passwordEntry.OnChanged = func(text string) {
		settings.Password = text
	}

	return container.NewVBox(
		usernameEntry,
		passwordEntry,
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
