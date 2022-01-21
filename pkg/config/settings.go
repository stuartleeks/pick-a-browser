package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/tidwall/jsonc"
)

type Settings struct {
	Browsers []Browser
	Rules    []Rule
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
	settingsFilename = path.Join(profilePath, settingsBaseFilename)
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

	// TODO - update check
	// TODO - transformations

	browsers, err := parseBrowsers(root)
	if err != nil {
		return nil, err
	}
	rules, err := parseRules(root)
	if err != nil {
		return nil, err
	}

	return &Settings{
		Browsers: browsers,
		Rules:    rules,
	}, nil
}

// func getObjectChildNode(node map[string]interface{}, propertyName string) (map[string]interface{}, error) {
// 	if childNode, ok := node[propertyName]; ok {
// 		switch childNode := childNode.(type) {
// 		case map[string]interface{}:
// 			return childNode, nil
// 		default:
// 			return map[string]interface{}{}, fmt.Errorf("%q property should be an object", propertyName)
// 		}
// 	}
// 	return map[string]interface{}{}, fmt.Errorf("%q property not found", propertyName)
// }

func getArrayChildNode(node map[string]interface{}, propertyName string, required bool) ([]interface{}, error) {
	if childNode, ok := node[propertyName]; ok {
		switch childNode := childNode.(type) {
		case []interface{}:
			return childNode, nil
		default:
			return []interface{}{}, fmt.Errorf("%q property should be an array", propertyName)
		}
	}
	if required {
		return []interface{}{}, fmt.Errorf("%q property not found", propertyName)
	}
	return []interface{}{}, nil
}

func getRequiredString(node map[string]interface{}, propertyName string) (string, error) {
	propertyNode, ok := node[propertyName]
	if !ok {
		return "", fmt.Errorf("required property %q not found", propertyName)
	}
	propertyValue, ok := propertyNode.(string)
	if !ok {
		return "", fmt.Errorf("required property %q expected to be a string", propertyName)
	}
	return propertyValue, nil
}

func getOptionalString(node map[string]interface{}, propertyName string) (*string, error) {
	propertyNode, ok := node[propertyName]
	if !ok {
		return nil, nil
	}
	propertyValue, ok := propertyNode.(string)
	if !ok {
		return nil, fmt.Errorf("optional property %q expected to be a string", propertyName)
	}
	return &propertyValue, nil
}

func getOptionalBoolWithDefault(node map[string]interface{}, propertyName string, defaultValue bool) (bool, error) {
	propertyNode, ok := node[propertyName]
	if !ok {
		return defaultValue, nil
	}
	propertyValue, ok := propertyNode.(bool)
	if !ok {
		return false, fmt.Errorf("optional property %q expected to be a bool", propertyName)
	}
	return propertyValue, nil
}
