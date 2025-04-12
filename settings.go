package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LauncherSettings struct {
	FortniteInstallPath string `json:"fortniteInstallPath"`
	Username            string `json:"username"`
	Password            string `json:"password"`
}

func loadSettings() LauncherSettings {
	localAppData := os.Getenv("LOCALAPPDATA")
	fortressAppData := filepath.Join(localAppData, ".FortressLauncher")
	settingsPath := filepath.Join(fortressAppData, "settings.json")

	jsonFile, err := os.Open(settingsPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var settings LauncherSettings
	json.Unmarshal(byteValue, &settings)

	return settings
}

func writeSettings(settings *LauncherSettings) {
	localAppData := os.Getenv("LOCALAPPDATA")
	fortressAppData := filepath.Join(localAppData, ".FortressLauncher")
	settingsPath := filepath.Join(fortressAppData, "settings.json")

	bytes, err := json.Marshal(settings)
	if err != nil {
		fmt.Println(err)
	}

	os.WriteFile(settingsPath, bytes, 0644)
}
