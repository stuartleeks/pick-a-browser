package config

import (
	"fmt"
	"log"
	"syscall"
)

type Browser struct {
	Id       string
	Name     string
	Exe      string
	Args     *string // TODO move to []string
	IconPath *string
	Hidden   bool
}

func (b *Browser) Launch(url string) error {
	// NOTE: I tried using os/exec.Command but hit issues with args :-(

	var sI syscall.StartupInfo
	var pI syscall.ProcessInformation

	exe := b.Exe
	if exe[0] != '"' {
		exe = "\"" + exe + "\""
	}
	runCommand := exe
	if b.Args != nil {
		runCommand += " " + *b.Args
	}
	if url != "" {
		runCommand += " " + url
	}

	log.Println(runCommand)
	argv := syscall.StringToUTF16Ptr(runCommand) //nolint:staticcheck
	err := syscall.CreateProcess(
		nil,
		argv,
		nil,
		nil,
		true,
		0,
		nil,
		nil,
		&sI,
		&pI)
	return err
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
