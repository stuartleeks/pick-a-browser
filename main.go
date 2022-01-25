//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lxn/walk"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"github.com/stuartleeks/pick-a-browser/pkg/appstate"
	"github.com/stuartleeks/pick-a-browser/pkg/config"
)

// Overridden via ldflags
var (
	version = "0.0.1-devbuild"
	commit  = "unknown"
	date    = "unknown"
)

const updateCheckInterval time.Duration = 4 * time.Hour // TODO - make configurable

func main() {
	settings, err := config.LoadSettings()
	if err != nil {
		walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Error loading settings:\n %s", err), walk.MsgBoxOK)
		return
	}

	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
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
			err := HandleUpdate()
			if err != nil {
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--browser-scan":
			err := HandleBrowserScan(settings)
			if err != nil {
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--version":
			walk.MsgBox(nil, "pick-a-browser", fmt.Sprintf("Version: %s (%s, %s)", version, commit, date), walk.MsgBoxOK)
			return
		}
	}

	if len(args) > 1 {
		walk.MsgBox(nil, "pick-a-browser error...", "Expected a single arg with the url", walk.MsgBoxOK)
		return
	}

	if err = PerformUpdateCheck(settings); err != nil {
		walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Failed to update:\n%s", err), walk.MsgBoxOK|walk.MsgBoxIconError)
	}

	url := ""
	if len(args) == 1 {
		url = args[0]
	}
	HandleUrl(url, settings)

}

func PerformUpdateCheck(settings *config.Settings) error {

	state, err := appstate.Load()
	if err != nil {
		return err
	}
	if settings.UpdateCheck == config.UpdateCheckAuto || settings.UpdateCheck == config.UpdateCheckPrompt {
		if time.Now().After(state.LastUpdateCheck.Add(updateCheckInterval)) {
			latest, err := CheckForUpdate(version)
			if err != nil {
				return err
			}
			state.LastUpdateCheck = time.Now().UTC()
			if err = appstate.Save(state); err != nil {
				return fmt.Errorf("error checking for updates:\n%s", err)
			}

			if latest == nil {
				return nil
			}

			// apply on auto
			apply := settings.UpdateCheck == config.UpdateCheckAuto
			if settings.UpdateCheck == config.UpdateCheckPrompt {
				result := walk.MsgBox(nil, "pick-a-browser update...", fmt.Sprintf("Version %s is available\n\n Update?", latest.Version), walk.MsgBoxYesNo|walk.MsgBoxIconQuestion)
				apply = result == walk.DlgCmdYes
			}

			if apply {
				exe, err := os.Executable()
				if err != nil {
					return fmt.Errorf("failed to locate executable:\n%s", err)
				}
				if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
					return fmt.Errorf("failed to perform update:\n%s", err)
				}

				walk.MsgBox(nil, "pick-a-browser...", fmt.Sprintf("Successfully updated to version %s", latest.Version), walk.MsgBoxOK)
			}
		}
	}
	return nil
}
