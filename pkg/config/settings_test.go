package config

import (
	"testing"

	"gotest.tools/assert"
)

func Test_ParseBrowsers(t *testing.T) {
	json := `{
		"browsers": [
			{
				"id" : "test1",
				"name" : "test1",
				"exe" : "test1_exe",
			},
			{
				"id" : "test2",
				"name" : "test2",
				"exe" : "test2_exe",
				"args" : "arg1 arg2",
				"iconPath" : "c:\\some\\path",
				"hidden" : true,
			}
		]
	}`

	settings, err := ParseSettings([]byte(json))
	assert.NilError(t, err)
	assert.Assert(t, settings != nil)

	assert.Assert(t, len(settings.Browsers) == 2)

	browser := settings.Browsers[0]
	assert.Equal(t, "test1", browser.Id)
	assert.Equal(t, "test1", browser.Name)
	assert.Equal(t, "test1_exe", browser.Exe)
	assert.Assert(t, browser.Args == nil)
	assert.Assert(t, browser.IconPath == nil)
	assert.Equal(t, false, browser.Hidden)

	browser = settings.Browsers[1]
	assert.Equal(t, "test2", browser.Id)
	assert.Equal(t, "test2", browser.Name)
	assert.Equal(t, "test2_exe", browser.Exe)
	assert.Assert(t, browser.Args != nil)
	assert.Equal(t, "arg1 arg2", *browser.Args)
	assert.Assert(t, browser.IconPath != nil)
	assert.Equal(t, "c:\\some\\path", *browser.IconPath)
	assert.Equal(t, true, browser.Hidden)
}

func Test_ParseLinkShorteners(t *testing.T) {
	json := `{
		"browsers": [],
		"transformations" : {
			"linkShorteners": [
				"shortener1",
				"shortener2",
			]
		},
		"rules": []
	}`

	settings, err := ParseSettings([]byte(json))
	assert.NilError(t, err)
	assert.Assert(t, settings != nil)

	assert.Assert(t, len(settings.Transformations.LinkShorteners) == 2)

	shortener := settings.Transformations.LinkShorteners[0]
	assert.Equal(t, "shortener1", shortener)

	shortener = settings.Transformations.LinkShorteners[1]
	assert.Equal(t, "shortener2", shortener)
}

func Test_ParseLinkWrappers(t *testing.T) {
	json := `{
		"browsers": [],
		"transformations" : {
			"linkWrappers": [
				{ "prefix": "https://example.com", "queryString" : "url" },
				{ "prefix": "https://example.net", "queryString" : "u" },
			]
		},
		"rules": []
	}`

	settings, err := ParseSettings([]byte(json))
	assert.NilError(t, err)
	assert.Assert(t, settings != nil)

	assert.Assert(t, len(settings.Transformations.LinkWrappers) == 2)

	wrapper := settings.Transformations.LinkWrappers[0]
	assert.Equal(t, "https://example.com", wrapper.UrlPrefix)
	assert.Equal(t, "url", wrapper.QueryStringKey)

	wrapper = settings.Transformations.LinkWrappers[1]
	assert.Equal(t, "https://example.net", wrapper.UrlPrefix)
	assert.Equal(t, "u", wrapper.QueryStringKey)
}

func Test_ParseRules(t *testing.T) {
	json := `{
		"browsers": [],
		"rules": [
			{
				"type" : "prefix",
				"prefix" : "https://example.com",
				"browser" : "browser1",
			},
			{
				"type" : "host",
				"host" : "example.com",
				"browser" : "browser2",
			},
		]
	}`

	settings, err := ParseSettings([]byte(json))
	assert.NilError(t, err)
	assert.Assert(t, settings != nil)

	assert.Assert(t, len(settings.Rules) == 2)

	rule := settings.Rules[0]
	assert.Assert(t, rule != nil)
	assert.Equal(t, "prefix", rule.Type())
	assert.Equal(t, "browser1", rule.BrowserId())

	prefixRule := rule.(*PrefixRule)
	assert.Equal(t, "https://example.com", prefixRule.prefixMatch)

	rule = settings.Rules[1]
	assert.Assert(t, rule != nil)
	assert.Equal(t, "host", rule.Type())
	assert.Equal(t, "browser2", rule.BrowserId())

	hostRule := rule.(*HostRule)
	assert.Equal(t, "example.com", hostRule.host)
}
