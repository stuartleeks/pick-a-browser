// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/lxn/walk"
	walkd "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows/registry"

	"github.com/stuartleeks/pick-a-browser/pkg/config"
)

func main() {
	settings, err := config.LoadSettings()
	if err != nil {
		walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Error loading settings:\n %s", err), walk.MsgBoxOK)
		return
	}

	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
		// TODO - handle --update, --install, --uninstall, ....
		case "--install":
			err := HandleInstall()
			if err == nil {
				walk.MsgBox(nil, "pick-a-browser", "Installed!", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			} else {
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--uninstall":
			err := HandleUninstall()
			if err == nil {
				walk.MsgBox(nil, "pick-a-browser", "Uninstalled!", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			} else {
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--update":
			walk.MsgBox(nil, "pick-a-browser TODO...", "--update not implemented yet", walk.MsgBoxOK|walk.MsgBoxIconError)
			return
		case "--browser-scan":
			walk.MsgBox(nil, "pick-a-browser TODO...", "--browser-scan not implemented yet", walk.MsgBoxOK|walk.MsgBoxIconError)
			return
		}
	}

	if len(args) > 1 {
		walk.MsgBox(nil, "pick-a-browser error...", "Expected a single arg with the url", walk.MsgBoxOK)
		return
	}

	url := ""
	if len(args) == 1 {
		url = args[0]
	}
	HandleUrl(url, settings)

}

type RegistryKey struct {
	BaseKey registry.Key
	Path    string
}
type RegistryValueToSet struct {
	BaseKey   registry.Key
	Path      string
	ValueName string
	Value     string
}

func HandleInstall() error {
	// Register as browser as per: https://docs.microsoft.com/en-us/windows/win32/shell/start-menu-reg

	startMenuInternetKey, _, err := registry.CreateKey(registry.LOCAL_MACHINE, "SOFTWARE\\Clients\\StartMenuInternet", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	_, openedExisting, err := registry.CreateKey(startMenuInternetKey, "pick-a-browser", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	if openedExisting {
		return fmt.Errorf("pick-a-browser is already installed")
	}

	exePath := os.Args[0]
	description := "browser selector - see https://github.com/stuartleeks/pick-a-browser"
	iconPath := fmt.Sprintf("%s,0", exePath)
	keysToCreate := []RegistryKey{
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\DefaultIcon"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\StartMenu"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\FileAssociations"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\UrlAssociations"},
		// InstallInfo: https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/cc144109(v=vs.85)
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\InstallInfo"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell\\open"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell\\open\\command"},
		// https://docs.microsoft.com/en-us/windows/win32/shell/default-programs#registering-an-application-for-use-with-default-programs
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\DefaultIcon"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell\\open"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell\\open\\command"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Application"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\StartMenu"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\FileAssociations"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\UrlAssociations"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\StartMenu"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\FileAssociation"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\UrlAssociations"},
	}
	valuesToSet := []RegistryValueToSet{
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser", ValueName: "", Value: "Pick A Browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\DefaultIcon", ValueName: "", Value: iconPath},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities", ValueName: "ApplicationName", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities", ValueName: "ApplicationDescription", Value: description},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\StartMenu", ValueName: "StartMenuInternet", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\UrlAssociations", ValueName: "http", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\UrlAssociations", ValueName: "https", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\InstallInfo", ValueName: "HideIconsCommand", Value: ""},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\InstallInfo", ValueName: "ReinstallCommand", Value: ""},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\InstallInfo", ValueName: "ShowIconsCommand", Value: ""},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\InstallInfo", ValueName: "IconsVisible", Value: "1"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell", ValueName: "", Value: "open"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell\\open\\command", ValueName: "", Value: exePath},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\DefaultIcon", ValueName: "", Value: iconPath},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell", ValueName: "", Value: "open"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell\\open\\command", ValueName: "", Value: fmt.Sprintf("\"%s\" \"%%1\"", exePath)},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Application", ValueName: "ApplicationName", Value: "pick-a-browser"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Application", ValueName: "ApplicationDescription", Value: description},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Application", ValueName: "ApplicationIcon", Value: description},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities", ValueName: "ApplicationName", Value: "pick-a-browser"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities", ValueName: "ApplicationDescription", Value: description},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\StartMenu", ValueName: "StartMenuInternet", Value: "pick-a-browser"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\UrlAssociations", ValueName: "http", Value: "pick-a-browser"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\UrlAssociations", ValueName: "https", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities", ValueName: "ApplicationName", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities", ValueName: "ApplicationDescription", Value: description},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\StartMenu", ValueName: "StartMenuInternet", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\UrlAssociations", ValueName: "http", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\UrlAssociations", ValueName: "https", Value: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\RegisteredApplications", ValueName: "pick-a-browser", Value: "SOFTWARE\\Clients\\StartmenuInternet\\pick-a-browser\\Capabilities"},
	}

	for _, keyToCreate := range keysToCreate {
		_, _, err := registry.CreateKey(keyToCreate.BaseKey, keyToCreate.Path, registry.ALL_ACCESS)
		if err != nil {
			return err
		}
	}

	for _, valueToSet := range valuesToSet {
		key, _, err := registry.CreateKey(valueToSet.BaseKey, valueToSet.Path, registry.ALL_ACCESS)
		if err != nil {
			return err
		}
		err = key.SetStringValue(valueToSet.ValueName, valueToSet.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleUninstall() error {

	startMenuInternetKey, _, err := registry.CreateKey(registry.LOCAL_MACHINE, "SOFTWARE\\Clients\\StartMenuInternet", registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("error opening 'SOFTWARE\\Clients\\StartMenuInternet': %s", err)
	}

	_, openedExisting, err := registry.CreateKey(startMenuInternetKey, "pick-a-browser", registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("error opening 'SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser': %s", err)
	}
	if !openedExisting {
		return fmt.Errorf("pick-a-browser is not installed")
	}

	keysToDelete := []RegistryKey{
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\UrlAssociations"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\FileAssociations"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities\\StartMenu"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\Capabilities"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\DefaultIcon"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\InstallInfo"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell\\open\\command"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell\\open"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser\\shell"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\Clients\\StartMenuInternet\\pick-a-browser"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\DefaultIcon"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell\\open\\command"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell\\open"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\shell"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Application"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\FileAssociations"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\UrlAssociations"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities\\StartMenu"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser\\Capabilities"},
		{BaseKey: registry.CLASSES_ROOT, Path: "pick-a-browser"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\StartMenu"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\FileAssociation"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities\\UrlAssociations"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser\\Capabilities"},
		{BaseKey: registry.LOCAL_MACHINE, Path: "SOFTWARE\\pick-a-browser"},
	}

	for _, keyToDelete := range keysToDelete {
		err := registry.DeleteKey(keyToDelete.BaseKey, keyToDelete.Path)
		if err != nil && !errors.Is(err, registry.ErrNotExist) {
			return fmt.Errorf("error opening %v %q: %s", getRootKeyName(keyToDelete.BaseKey), keyToDelete.Path, err)
		}
	}

	registeredAppliationsKey, _, err := registry.CreateKey(registry.LOCAL_MACHINE, "SOFTWARE\\RegisteredApplications", registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("error opening 'SOFTWARE\\RegisteredApplications': %s", err)
	}
	err = registeredAppliationsKey.DeleteValue("pick-a-browser")
	if err != nil {
		return fmt.Errorf("error deleting 'HKLM\\SOFTWARE\\RegisteredApplications!pick-a-browser': %s", err)
	}

	return nil
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

func HandleUrl(urlString string, settings *config.Settings) {
	if urlString != "" {
		url, err := url.Parse(urlString)
		if err != nil {
			walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Failed to parse url %q:\n%s", urlString, err), walk.MsgBoxOK)
			return
		}

		linkWrappers := append(config.GetDefaultLinkWrappers(), settings.Transformations.LinkWrappers...)
		linkShorteners := append(config.GetDefaultLinkShorteners(), settings.Transformations.LinkShorteners...)
		for {
			// TODO - decide what to do with errors here
			newUrl, _ := transformUrlWithWrappers(url, linkWrappers)
			if newUrl != nil {
				url = newUrl
				continue
			}

			newUrl, _ = transformUrlWithShorteners(url, linkShorteners)
			if newUrl != nil {
				url = newUrl
				continue
			}

			break
		}
		urlString = url.String()

		// match rules - launch browser and exit on match, or fall through to show list
		matchedBrowserId := matchRules(settings.Rules, url)

		if matchedBrowserId != "" {
			// get browser from Id
			for _, browser := range settings.Browsers {
				if browser.Id == matchedBrowserId {
					err := browser.Launch(urlString)
					if err != nil {
						walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Failed to launch browser (%q):\n%s", matchedBrowserId, err), walk.MsgBoxOK)
					}
					return
				}
			}
			walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Failed to find browser with id %q", matchedBrowserId), walk.MsgBoxOK)
			return
		}
	}

	// If here then show browser list (filter to non-hidden browsers)
	browsers := []config.Browser{}
	for _, browser := range settings.Browsers {
		if !browser.Hidden {
			browsers = append(browsers, browser)
		}
	}

	mw := &MyMainWindow{}

	defaultFont := walkd.Font{Family: "Segoe UI", PointSize: 16}

	widgets := []walkd.Widget{
		walkd.Label{
			AssignTo: &mw.urlLabel,
			Font:     defaultFont,
			Text:     "URL: " + urlString,
		},
	}

	for _i, tmp := range browsers {
		browserNumber := _i + 1
		browser := tmp
		widgets = append(widgets, walkd.PushButton{
			Text:      fmt.Sprintf("&%d: %s", browserNumber, browser.Name),
			Row:       browserNumber,
			MinSize:   walkd.Size{Width: 150, Height: 10},
			MaxSize:   walkd.Size{Width: 3000, Height: 150},
			Font:      walkd.Font{Family: "Segoe UI", PointSize: 20},
			Alignment: walkd.AlignHCenterVNear,
			OnClicked: func() {
				// TODO launch!
				err := browser.Launch(urlString)
				if err != nil {
					walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Failed to launch browser (%q):\n%s", browser.Id, err), walk.MsgBoxOK)
				}
				mw.Close()
			},
		})
	}

	window := walkd.MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "pick-a-browser...",
		MinSize:  walkd.Size{Width: 150, Height: 150},
		Size:     walkd.Size{Width: 300, Height: 400},
		// Layout:   walkd.VBox{MarginsZero: true},
		Layout: walkd.Grid{
			MarginsZero: true,
			Rows:        len(settings.Browsers) + 1,
			Margins:     walkd.Margins{Top: 15, Bottom: 15, Left: 0, Right: 0},
			// SpacingZero: true,
			// Spacing:     0,
			Alignment: walkd.AlignHCenterVNear,
		},
		Children: widgets,
		// OnKeyDown isn't being invoked :-(
		// OnKeyDown: func(key walk.Key) {
		// 	log.Println("keydown")
		// 	walk.MsgBox(mw, "test", fmt.Sprintf("key: %v", key), walk.MsgBoxOK)
		// 	browserIndex := -1
		// 	switch key {
		// 	case walk.Key1:
		// 		browserIndex = 0
		// 	case walk.Key2:
		// 		browserIndex = 1
		// 	case walk.Key3:
		// 		browserIndex = 2
		// 	case walk.Key4:
		// 		browserIndex = 3
		// 	case walk.Key5:
		// 		browserIndex = 4
		// 	case walk.Key6:
		// 		browserIndex = 5
		// 	case walk.Key7:
		// 		browserIndex = 6
		// 	case walk.Key8:
		// 		browserIndex = 7
		// 	case walk.Key9:
		// 		browserIndex = 8
		// 	default:
		// 		return
		// 	}
		// 	browser := settings.Browsers[browserIndex]
		// 	walk.MsgBox(mw, "test", fmt.Sprintf("TODO: launch %q", browser.Name), walk.MsgBoxOK)
		// },
	}

	if _, err := window.Run(); err != nil {
		log.Fatal(err)
	}

}

type MyMainWindow struct {
	*walk.MainWindow
	urlLabel *walk.Label
}

// TODO move out of main and add tests
func matchRules(rules []config.Rule, url *url.URL) string {
	matchWeight := -1
	browserId := ""
	for _, rule := range rules {
		tmpWeight := rule.Match(url)
		if tmpWeight > matchWeight {
			matchWeight = tmpWeight
			browserId = rule.BrowserId()
		}
	}
	return browserId
}

func transformUrlWithShorteners(url *url.URL, linkShorteners []string) (*url.URL, error) {
	for _, linkShortener := range linkShorteners {
		if strings.HasSuffix(url.Host, linkShortener) {
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Get(url.String())
			if err != nil {
				return nil, err
			}
			newUrlString := resp.Header.Get("Location")
			return url.Parse(newUrlString)
		}
	}
	return nil, nil
}
func transformUrlWithWrappers(url *url.URL, linkWrappers []config.LinkWrapper) (*url.URL, error) {
	for _, linkWrapper := range linkWrappers {
		if strings.HasPrefix(url.String(), linkWrapper.UrlPrefix) {
			newUrlString := url.Query().Get(linkWrapper.QueryStringKey)
			if newUrlString == "" {
				return nil, nil
			}
			return url.Parse(newUrlString)
		}
	}
	return nil, nil
}
