package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/jsonc"
)

type UpdateCheck string

const (
	UpdateCheckNone   UpdateCheck = "never"
	UpdateCheckPrompt UpdateCheck = "prompt"
	UpdateCheckAuto   UpdateCheck = "auto"
)

type Settings struct {
	Browsers        []Browser
	Transformations Transformations
	Rules           []Rule
	UpdateCheck     UpdateCheck
}

const settingsBaseFilename string = "pick-a-browser-settings.json"

func getSettingsFilename() (string, error) {
	// If PICK_A_BROWSER_CONFIG is set, use it
	settingsFilename := os.Getenv("PICK_A_BROWSER_CONFIG")
	if settingsFilename != "" {
		return settingsFilename, nil
	}

	// Try user profile settings file first if it exists...
	profilePath := os.Getenv("USERPROFILE")
	if profilePath == "" {
		return "", fmt.Errorf("USERPROFILE env var not set")
	}
	settingsFilename = filepath.Join(profilePath, settingsBaseFilename)
	_, err := os.Stat(settingsFilename)
	if err != nil {
		return settingsFilename, nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("error accessing settings file (%q): %s", settingsFilename, err)
	}

	// Lastly, look for settings next to the app
	exe := os.Args[0]
	return filepath.Join(filepath.Dir(exe), settingsBaseFilename), nil
}

func LoadSettings() (*Settings, error) {
	settingsFilename, err := getSettingsFilename()
	if err != nil {
		return nil, err
	}
	return LoadSettingsFromFile(settingsFilename)
}
func LoadSettingsFromFile(filename string) (*Settings, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseSettings(buf)
}
func ParseSettings(jsonBuf []byte) (*Settings, error) {
	var root map[string]interface{}
	jsonBuf = jsonc.ToJSON(jsonBuf)
	err := json.Unmarshal(jsonBuf, &root)
	if err != nil {
		return nil, err
	}

	browsers, err := parseBrowsers(root)
	if err != nil {
		return nil, err
	}
	transformations, err := parseTransformations(root)
	if err != nil {
		return nil, err
	}
	rules, err := parseRules(root)
	if err != nil {
		return nil, err
	}
	updateCheck, err := parseUpdateCheck(root)
	if err != nil {
		return nil, err
	}

	return &Settings{
		Browsers:        browsers,
		Transformations: transformations,
		Rules:           rules,
		UpdateCheck:     updateCheck,
	}, nil
}

func parseUpdateCheck(rootNode map[string]interface{}) (UpdateCheck, error) {
	updateCheckString, err := getOptionalString(rootNode, "updates")
	if err != nil {
		return "", err
	}

	if updateCheckString == nil {
		return UpdateCheckAuto, nil
	}

	switch UpdateCheck(*updateCheckString) {
	case UpdateCheckNone, UpdateCheckPrompt, UpdateCheckAuto:
		return UpdateCheck(*updateCheckString), nil
	default:
		return "", fmt.Errorf("unrecognised value for updateCheck (%q), supported values are none, prompt, auto", *updateCheckString)
	}
}
