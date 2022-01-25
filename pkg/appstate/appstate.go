package appstate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type AppState struct {
	LastUpdateCheck time.Time `json:"lastUpdateCheck"`
}

func Load() (AppState, error) {
	appDataPath := os.Getenv("LOCALAPPDATA")

	filename := filepath.Join(appDataPath, "stuartleeks", "pick-a-browser", "app-state.json")

	buf, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// no file - return default
			return AppState{}, nil
		}
		return AppState{}, err
	}

	var appState AppState

	if err := json.Unmarshal(buf, &appState); err != nil {
		return AppState{}, err
	}

	return appState, nil
}
func Save(appState AppState) error {
	buf, err := json.Marshal(appState)
	if err != nil {
		return err
	}

	appDataPath := os.Getenv("LOCALAPPDATA")
	filename := filepath.Join(appDataPath, "stuartleeks", "pick-a-browser", "app-state.json")
	statePath, _ := filepath.Split(filename)

	if err = os.MkdirAll(statePath, 0666); err != nil {
		return err
	}
	if err = os.WriteFile(filename, buf, 0666); err != nil {
		return err
	}
	return nil
}
