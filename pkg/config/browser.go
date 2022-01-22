package config

import (
	"fmt"
)

type Browser struct {
	Id       string
	Name     string
	Exe      string
	Args     *string // TODO move to []string
	IconPath *string
	Hidden   bool
}

// TODO - browser scan

func parseBrowsers(rootNode map[string]interface{}) ([]Browser, error) {
	browsersNode, err := getArrayChildNode(rootNode, "browsers", true)
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
