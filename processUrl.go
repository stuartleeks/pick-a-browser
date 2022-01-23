// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/lxn/walk"
	walkd "github.com/lxn/walk/declarative"

	"github.com/stuartleeks/pick-a-browser/pkg/config"
)

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
