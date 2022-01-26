//go:build windows
// +build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

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
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("failed to convert to absolute path (%q): %s", exePath, err)
	}

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
