//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lxn/walk"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	log "github.com/sirupsen/logrus"
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

	err = configureLogger(settings)
	if err != nil {
		walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Error creating logger:\n %s", err), walk.MsgBoxOK)
	}
	logger := log.WithField("PID", os.Getpid())

	args := os.Args[1:]
	log.Infoln("pick-a-browser.exe starting ", args)

	if len(args) > 0 {
		switch args[0] {
		case "--install":
			logger.Infoln("Installing")
			err := HandleInstall()
			if err == nil {
				logger.Infoln("Install succeeded")
				walk.MsgBox(nil, "pick-a-browser", "Installed!", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			} else {
				logger.Errorln("Install failed", err)
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--uninstall":
			logger.Infoln("Uninstalling")
			err := HandleUninstall()
			if err == nil {
				logger.Infoln("Uninstall succeeded")
				walk.MsgBox(nil, "pick-a-browser", "Uninstalled!", walk.MsgBoxOK|walk.MsgBoxIconInformation)
			} else {
				logger.Errorln("Uninstall failed", err)
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--update":
			logger.Infoln("Updating")
			err := HandleUpdate()
			if err == nil {
				logger.Infoln("Update failed", err)
			} else {
				logger.Errorln("Update failed", err)
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--browser-scan":
			logger.Infoln("BrowserScan")
			err := HandleBrowserScan(settings)
			if err == nil {
				logger.Infoln("BrowserScan completed")
			} else {
				logger.Errorln("BrowserScan failed", err)
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
			return
		case "--version":
			logger.Traceln("Showing version dialog")
			walk.MsgBox(nil, "pick-a-browser", fmt.Sprintf("Version: %s (%s, %s)", version, commit, date), walk.MsgBoxOK)
			return
		}
	}

	if len(args) > 1 {
		logger.Errorln("Invalid args", args)
		walk.MsgBox(nil, "pick-a-browser error...", "Expected a single arg with the url", walk.MsgBoxOK)
		return
	}

	if err = PerformUpdateCheck(settings); err != nil {
		logger.Errorln("UpdateCheck failed", err)
		walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Failed to update:\n%s", err), walk.MsgBoxOK|walk.MsgBoxIconError)
	}

	url := ""
	if len(args) == 1 {
		url = args[0]
	}
	err = HandleUrl(url, settings)
	if err != nil {
		logger.Errorln("HandleUrl failed", err)
		walk.MsgBox(nil, "pick-a-browser error...", fmt.Sprintf("Error handling URL:\n%s", err), walk.MsgBoxOK|walk.MsgBoxIconError)
	}
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
				if err := selfupdate.NoGitUpdater().UpdateTo(latest, exe); err != nil {
					return fmt.Errorf("failed to perform update:\n%s", err)
				}

				walk.MsgBox(nil, "pick-a-browser...", fmt.Sprintf("Successfully updated to version %s", latest.Version), walk.MsgBoxOK)
			}
		}
	}
	return nil
}

func configureLogger(settings *config.Settings) error {
	appDataPath := os.Getenv("LOCALAPPDATA")
	oldFilename := filepath.Join(appDataPath, "stuartleeks", "pick-a-browser", "pick-a-browser-old.log")
	filename := filepath.Join(appDataPath, "stuartleeks", "pick-a-browser", "pick-a-browser.log")
	statePath, _ := filepath.Split(filename)

	if err := os.MkdirAll(statePath, 0666); err != nil {
		return err
	}

	const maxFileSize = 1024 * 1024 // Trigger log rotation at 1MB
	if info, err := os.Stat(filename); err == nil {
		if info.Size() > maxFileSize {
			if _, err = os.Stat(oldFilename); err == nil {
				if err = os.Remove(oldFilename); err != nil {
					return fmt.Errorf("failed to delete old log file: %s", err)
				}
			}
			if err = os.Rename(filename, oldFilename); err != nil {
				return fmt.Errorf("failed to rename log file: %s", err)
			}
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file %q: %s", filename, err)
	}
	log.SetOutput(file)
	log.SetLevel(getLogLevel(settings.Log.Level))
	log.SetFormatter(&log.TextFormatter{})

	return nil
}

func getLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.ErrorLevel
	}
}
