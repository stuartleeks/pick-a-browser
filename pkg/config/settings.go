package config

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/jsonc"
)

type Settings struct {
	Browsers []Browser
}

type Browser struct {
	Id       string
	Name     string
	Exe      string
	Args     *string
	IconPath *string
	Hidden   bool
}

func ParseSettingsFromJson(jsonString string) (*Settings, error) {

	var root map[string]interface{}
	jsonBuf := jsonc.ToJSON([]byte(jsonString))
	err := json.Unmarshal(jsonBuf, &root)
	if err != nil {
		return nil, err
	}

	// TODO - can we handle JSONC??

	// TODO - update check
	// TODO - transformations
	// TODO - rules

	browsers, err := parseBrowsers(root)
	if err != nil {
		return nil, err
	}

	return &Settings{
		Browsers: browsers,
	}, nil
}

// TODO - browser scan
// TODO - browser launch

func parseBrowsers(rootNode map[string]interface{}) ([]Browser, error) {
	browsersNode, err := getArrayChildNode(rootNode, "browsers")
	if err != nil {
		return []Browser{}, err
	}

	browsers := []Browser{}

	for _, v := range browsersNode {
		browserNode, ok := v.(map[string]interface{})
		if !ok {
			return []Browser{}, fmt.Errorf("expected objects in browsers array")
		}
		browser, err := parseBrowser(browserNode)
		if err != nil {
			return []Browser{}, err
		}
		browsers = append(browsers, browser)
	}

	return browsers, nil
}
func parseBrowser(browserNode map[string]interface{}) (Browser, error) {

	id, err := getRequiredString(browserNode, "id")
	if err != nil {
		return Browser{}, err
	}
	name, err := getRequiredString(browserNode, "name")
	if err != nil {
		return Browser{}, err
	}
	exe, err := getRequiredString(browserNode, "exe")
	if err != nil {
		return Browser{}, err
	}
	args, err := getOptionalString(browserNode, "args")
	if err != nil {
		return Browser{}, err
	}
	iconPath, err := getOptionalString(browserNode, "iconPath")
	if err != nil {
		return Browser{}, err
	}
	hidden, err := getOptionalBoolWithDefault(browserNode, "hidden", false)
	if err != nil {
		return Browser{}, err
	}

	return Browser{
		Id:       id,
		Name:     name,
		Exe:      exe,
		Args:     args,
		IconPath: iconPath,
		Hidden:   hidden,
	}, nil
}

func getObjectChildNode(node map[string]interface{}, propertyName string) (map[string]interface{}, error) {
	if childNode, ok := node[propertyName]; ok {
		switch childNode := childNode.(type) {
		case map[string]interface{}:
			return childNode, nil
		default:
			return map[string]interface{}{}, fmt.Errorf("%q property should be an object", propertyName)
		}
	}
	return map[string]interface{}{}, fmt.Errorf("%q property not found", propertyName)
}

func getArrayChildNode(node map[string]interface{}, propertyName string) ([]interface{}, error) {
	if childNode, ok := node[propertyName]; ok {
		switch childNode := childNode.(type) {
		case []interface{}:
			return childNode, nil
		default:
			return []interface{}{}, fmt.Errorf("%q property should be an array", propertyName)
		}
	}
	return []interface{}{}, fmt.Errorf("%q property not found", propertyName)
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
