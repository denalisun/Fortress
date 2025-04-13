package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.RGBA{
			R: 20,
			G: 20,
			B: 40,
			A: 0,
		}
	} else if name == theme.ColorNameButton {
		return color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 17
	}

	return theme.DefaultTheme().Size(name)
}

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
