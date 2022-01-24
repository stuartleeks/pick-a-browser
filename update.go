//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blang/semver"
	"github.com/lxn/walk"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"github.com/stuartleeks/pick-a-browser/pkg/appstate"
)

func HandleUpdate() error {

	latest, err := CheckForUpdate(version)
	if err != nil {
		return err
	}

	state, err := appstate.Load()
	if err != nil {
		return err
	}
	state.LastUpdateCheck = time.Now().UTC()
	if err = appstate.Save(state); err != nil {
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
