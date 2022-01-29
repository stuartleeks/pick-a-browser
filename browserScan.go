//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	walkd "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows/registry"

	"github.com/stuartleeks/pick-a-browser/pkg/config"
)

func HandleBrowserScan(settings *config.Settings) error {

	allBrowsers, err := getAllBrowsers()
	if err != nil {
		return err
	}

	mergedBrowsers := mergeSettings(settings.Browsers, allBrowsers)

	newSettings := map[string]interface{}{
		"browsers": mergedBrowsers,
	}

	buf, err := json.MarshalIndent(newSettings, "", "\t")
	if err != nil {
		return err
	}

	settingsString := strings.ReplaceAll(string(buf), "\n", "\r\n")

	window := walkd.MainWindow{
		Title:   "pick-a-browser...",
		MinSize: walkd.Size{Width: 500, Height: 400},
		Size:    walkd.Size{Width: 500, Height: 400},
		Layout:  walkd.VBox{MarginsZero: true},
		Children: []walkd.Widget{
			walkd.TextEdit{
				Text:               settingsString,
				ReadOnly:           true,
				AlwaysConsumeSpace: true,
				VScroll:            true,
			},
		},
	}

	if _, err := window.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func mergeSettings(baseBrowsers []config.Browser, newBrowsers []config.Browser) []config.Browser {
	result := []config.Browser{}
	result = append(result, baseBrowsers...)

	for _, newBrowser := range newBrowsers {
		found := hasMatchOnExeAndArgs(baseBrowsers, newBrowser)
		if !found {
			result = append(result, newBrowser)
		}
	}
	return result
}

func hasMatchOnExeAndArgs(browsers []config.Browser, browserToMatch config.Browser) bool {
	for _, browser := range browsers {
		if browser.Exe == browserToMatch.Exe &&
			((browser.Args == nil && browserToMatch.Args == nil) ||
				*browser.Args == *browserToMatch.Args) {
			return true
		}
	}
	return false
}

func getAllBrowsers() ([]config.Browser, error) {
	hklmBrowsers, err := getBrowsersFor(registry.LOCAL_MACHINE)
	if err != nil {
		return []config.Browser{}, fmt.Errorf("failed to local HKLM browsers: %s", err)
	}
	hkcuBrowsers, err := getBrowsersFor(registry.CURRENT_USER)
	if err != nil {
		return []config.Browser{}, fmt.Errorf("failed to local HKCU browsers: %s", err)
	}

	browsersTemp := append(hkcuBrowsers, hklmBrowsers...)
	browsersExpanded := []config.Browser{}

	for _, browser := range browsersTemp {
		switch browser.Name {
		case "Microsoft Edge":
			edgeProfiles, err := getEdgeProfiles()
			if err != nil {
				return []config.Browser{}, fmt.Errorf("failed to get Edge profiles: %s", err)
			}
			for _, edgeProfile := range edgeProfiles {
				args := fmt.Sprintf("--profile-directory=\"%s\"", edgeProfile)
				browsersExpanded = append(browsersExpanded, config.Browser{
					Id:       uuid.NewString(),
					Name:     browser.Name + " - " + edgeProfile,
					Exe:      browser.Exe,
					Args:     &args,
					IconPath: browser.IconPath,
					Hidden:   false,
				})
			}
		case "Google Chrome":
			edgeProfiles, err := getChromeProfiles()
			if err != nil {
				return []config.Browser{}, fmt.Errorf("failed to get Chrome profiles: %s", err)
			}
			for _, edgeProfile := range edgeProfiles {
				args := fmt.Sprintf("--profile-directory=\"%s\"", edgeProfile)
				browsersExpanded = append(browsersExpanded, config.Browser{
					Id:       uuid.NewString(),
					Name:     browser.Name + " - " + edgeProfile,
					Exe:      browser.Exe,
					Args:     &args,
					IconPath: browser.IconPath,
					Hidden:   false,
				})
			}
		default:
			browsersExpanded = append(browsersExpanded, browser)
		}
	}
	return browsersExpanded, nil
}

const reg_READ = registry.ENUMERATE_SUB_KEYS | registry.QUERY_VALUE

func getBrowsersFor(rootKey registry.Key) ([]config.Browser, error) {

	browsersKey, err := registry.OpenKey(rootKey, "SOFTWARE\\Clients\\StartMenuInternet", reg_READ)
	if err != nil {
		return []config.Browser{}, fmt.Errorf("failed to open SOFTWARE\\Clients\\StartMenuInternet")
	}

	browserKeyNames, err := browsersKey.ReadSubKeyNames(-1)
	if err != nil {
		return []config.Browser{}, fmt.Errorf("failed to get subkeys for SOFTWARE\\Clients\\StartMenuInternet")
	}

	browsers := []config.Browser{}
	for _, browserKeyName := range browserKeyNames {
		if browserKeyName == "pick-a-browser" {
			continue
		}
		browserKey, err := registry.OpenKey(browsersKey, browserKeyName, reg_READ)
		if err != nil {
			return []config.Browser{}, fmt.Errorf("failed to get open browser key %q", browserKeyName)
		}
		browser, err := getBrowserFromRegistry(browserKey, browserKeyName)
		// TODO - decide whether to skip on error or fail?
		if err != nil {
			return []config.Browser{}, fmt.Errorf("failed to parse browser settings (%q): %s", browserKeyName, err)
		}
		browsers = append(browsers, browser)
	}

	return browsers, nil
}

func getBrowserFromRegistry(key registry.Key, keyName string) (config.Browser, error) {

	commandKey, err := registry.OpenKey(key, "shell\\open\\command", reg_READ)
	if err != nil {
		return config.Browser{}, fmt.Errorf("failed to open shell\\open\\command key: %s", err)
	}

	exe, _, err := commandKey.GetStringValue("")
	if err != nil {
		return config.Browser{}, fmt.Errorf("failed to get value for shell\\open\\command key: %s", err)
	}
	exe = strings.Trim(exe, "\"")

	name, _, err := key.GetStringValue("")
	if err != nil {
		name = keyName
	}

	var iconPath *string
	iconKey, err := registry.OpenKey(key, "DefaultIcon", reg_READ)
	if err == nil {
		iconValue, _, err := iconKey.GetStringValue("")
		if err != nil {
			return config.Browser{}, fmt.Errorf("failed to get DefaultIcon value: %s", err)
		}
		iconPath = &iconValue
	} else {
		if !errors.Is(err, registry.ErrNotExist) {
			return config.Browser{}, fmt.Errorf("failed to open DefaultIcon key: %s", err)
		}
	}

	return config.Browser{
		Id:       uuid.NewString(),
		Name:     name,
		Exe:      exe,
		IconPath: iconPath,
		Hidden:   false,
	}, nil

}

func getEdgeProfiles() ([]string, error) {
	userProfile := os.Getenv("LOCALAPPDATA")
	edgeUserDataPath := filepath.Join(userProfile, "Microsoft\\Edge\\User Data")

	entries, err := os.ReadDir(edgeUserDataPath)
	if err != nil {
		return []string{}, nil
	}
	profiles := []string{"Default"}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "Profile ") {
			profiles = append(profiles, entry.Name())
		}
	}
	return profiles, nil
}

func getChromeProfiles() ([]string, error) {
	userProfile := os.Getenv("LOCALAPPDATA")
	chromeUserDataPath := filepath.Join(userProfile, "Google\\Chrome\\User Data")

	entries, err := os.ReadDir(chromeUserDataPath)
	if err != nil {
		return []string{}, nil
	}
	profiles := []string{"Default"}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "Profile ") {
			profiles = append(profiles, entry.Name())
		}
	}
	return profiles, nil
}

func getRootKeyName(key registry.Key) string {
	switch key {
	case registry.CLASSES_ROOT:
		return "HKCR"
	case registry.LOCAL_MACHINE:
		return "HKLM"
	default:
		return "unhandled root key"
	}
}
