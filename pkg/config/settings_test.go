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

	settings, err := ParseSettingsFromJson(json)
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
