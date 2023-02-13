package config

import (
	"fmt"
)

// TODO - look at whether to collapse the json parsing down to the standard library. Want to be able to handle args being single string vs array etc...
//
//	Also look at custom marhsalling: http://choly.ca/post/go-json-marshalling/
type Browser struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Exe      string  `json:"exe"`
	Args     *string `json:"args"` // TODO move to []string to allow prompting for limited subset when matched
	IconPath *string `json:"iconPath"`
	Hidden   bool    `json:"hidden"`
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
