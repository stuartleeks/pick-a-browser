package config

import (
	"net/url"
	"testing"

	"gotest.tools/assert"
)

func urlMustParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
func Test_SimplePrefixMatch(t *testing.T) {
	json := `{
		"browsers": [],
		"rules": [{
            "type": "prefix",
            "prefix": "https://example.com",
            "browser": "browser1"
        },{
            "type": "prefix",
            "prefix": "https://example.com/somewhere",
            "browser": "browser2"
        }
		]
	}`

	settings, err := ParseSettings([]byte(json))
	assert.NilError(t, err)
	assert.Assert(t, settings != nil)

	assert.Assert(t, len(settings.Rules) == 2)

	assert.Equal(t, MatchRules(settings.Rules, urlMustParse("https://example.com")), "browser1")
}

func Test_SimplePrefixMatchesLongest(t *testing.T) {
	json := `{
		"browsers": [],
		"rules": [{
            "type": "prefix",
            "prefix": "https://example.com",
            "browser": "browser1"
        },{
            "type": "prefix",
            "prefix": "https://example.com/somewhere",
            "browser": "browser2"
        }
		]
	}`

	settings, err := ParseSettings([]byte(json))
	assert.NilError(t, err)
	assert.Assert(t, settings != nil)

	assert.Assert(t, len(settings.Rules) == 2)

	assert.Equal(t, MatchRules(settings.Rules, urlMustParse("https://example.com/somewhere")), "browser2")
}
