// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/lxn/walk"

	"github.com/stuartleeks/pick-a-browser/pkg/config"
)

// Overridden via ldflags
var (
	version   = "0.0.1-devbuild"
	commit    = "unknown"
	date      = "unknown"
	goversion = "unknown"
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
			walk.MsgBox(nil, "pick-a-browser TODO...", "--update not implemented yet", walk.MsgBoxOK|walk.MsgBoxIconError) // TODO
			return
		case "--browser-scan":
			err := HandleBrowserScan(settings)
			if err != nil {
				walk.MsgBox(nil, "pick-a-browser", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
			}
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
