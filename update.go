//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/lxn/walk"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func HandleUpdate() error {

	latest, err := CheckForUpdate(version)
	if err != nil {
		return err
	}

	if latest == nil {
		walk.MsgBox(nil, "pick-a-browser...", fmt.Sprintf("Already on the latest version (%s)", version), walk.MsgBoxOK)
		return nil
	}

	result := walk.MsgBox(nil, "pick-a-browser...", fmt.Sprintf("Update to %s?\n\n%s", latest.Version, latest.ReleaseNotes), walk.MsgBoxYesNo)

	if result == walk.DlgCmdYes {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("could not locate executable path: %v", err)
		}
		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			return fmt.Errorf("error occurred while updating binary: %v", err)
		}
		fmt.Printf("Successfully updated to version %s\n", latest.Version)
	}
	return nil
}

func CheckForUpdate(currentVersion string) (*selfupdate.Release, error) {

	latest, found, err := selfupdate.DetectLatest("stuartleeks/pick-a-browser")
	if err != nil {
		return nil, fmt.Errorf("error occurred while detecting version: %v", err)
	}

	v, err := semver.Parse(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("error occurred while parsing version: %v", err)
	}

	if !found || latest.Version.LTE(v) {
		return nil, nil
	}
	return latest, nil
}

// func PeriodicCheckForUpdate(currentVersion string) {
// 	const checkInterval time.Duration = 24 * time.Hour

// 	lastCheck := config.GetLastUpdateCheck()

// 	if time.Now().Before(lastCheck.Add(checkInterval)) {
// 		return
// 	}
// 	fmt.Println("Checking for updates...")
// 	latest, err := CheckForUpdate(currentVersion)
// 	if err != nil {
// 		fmt.Printf("Error checking for updates: %s", err)
// 	}

// 	config.SetLastUpdateCheck(time.Now())
// 	if err = config.SaveConfig(); err != nil {
// 		fmt.Printf("Error saving last update check time: :%s\n", err)
// 	}

// 	if latest == nil {
// 		return
// 	}

// 	fmt.Printf("\n\n UPDATE AVAILABLE: %s \n \n Release notes: %s\n", latest.Version, latest.ReleaseNotes)
// 	fmt.Printf("Run `devcontainer update` to apply the update\n\n")
// }
