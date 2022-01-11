# pick-a-browser

`pick-a-browser` is a browser selector for Windows.
It registers as a browser and when launched it runs a configurable rule set to determine the real browser to launch.
If no rule matches then it prompts for the browser to launch from a configured list of browsers/profiles.

`pick-a-browser` was born after happily using (and contributing to) [BrowserPicker](https://github.com/mortenn/BrowserPicker).
After a number of reinstallations of the operating system, I wanted an automatic way to save my browser and rules configuration.
As a result, I created `pick-a-browser` as a way to address this (and to tweak behaviour in ways that suit my usage patterns slightly better).

## Installing

-   Get the binaries
-   Install `pick-a-browser` as a browser
-   Create `pick-a-browser-settings.json`
-   Set `pick-a-browser` as your default browser - <https://support.microsoft.com/en-gb/windows/change-your-default-browser-in-windows-10-020c58c6-7d77-797a-b74e-8f07946c5db6>

### Get the binaries

Currently, you can either grab the build artifact from [the latest CI build](https://github.com/stuartleeks/pick-a-browser/actions/workflows/ci-build.yml) and unzip, or clone and build from source.

### Installing pick-a-browser

To install `pick-a-browser`, run `pick-a-browser --install` (needs elevated permissions). You can use `pick-a-browser --uninstall` to uninstall

### Create settings file

See [Configuration](#configuration) below.

## Configuration

By default, `pick-a-browser` will look for `pick-a-browser-settings.json` in your user profile folder and then in the same folder as the app itself.

If you which to put the settings in a different location, set the `PICK_A_BROWSER_CONFIG` environment variable to the full path to the settings file.

### Browsers

The top-level `browsers` property allows you to configure browsers (or browser profiles) that `pick-a-browser` should use.

You can run `pick-a-browser --browser-scan` to generate the initial `browsers` section of the configuration and then copy and paste this into the config file.

The `browsers` property is an array of objects with the following properties:

| Name     | Type              | Description                                                                  |
| -------- | ----------------- | ---------------------------------------------------------------------------- |
| id       | string (required) | The id to use to identify the browser - used to refer to browsers in rules   |
| name     | string (required) | The name to display in the browser picker UI                                 |
| exe      | string (required) | The path to the app to launch for the browser                                |
| args     | string (optional) | Any arguments to pass to the browser. Useful for specifying browser profiles |
| iconPath | string (optional) | not currently used                                                           |
| hidden   | bool (optional)   | When no rules are matched, the UI displays a list of all non-hidden browsers |

### Transformations

The top-level `transformations` property allows you to define URL transformations that are applied before rules are evaluated.

There is some built-in handling of link shorteners/wrappers, but you can

#### Link Shorteners

Link shorteners provide a short link that redirects to the underlying link.

To add a link shortener, add the host value to the `linkShorteners` array under the `transformations` property.

```jsonc
{
	"transformations" : {
		"linkShorteners": [ "aka.ms" ] // NOTE that `aka.ms` is included in the default shorteners
	}
}
```

#### Link Wrappers

Link wrappers provide a URL that embeds the underlying URL in a query string parameter.

To add a link wrapper, add the wrapper value to the `linkWrappers` array under the `transformations` property.

| Name        | Type              | Description                                              |
| ----------- | ----------------- | -------------------------------------------------------- |
| prefix      | string (required) | The URL prefix to match for the shortener                |
| queryString | string (required) | The name of the query string value that contains the URL |

```jsonc
{
	"transformations" : {
		"linkWrappers" : [
			{ "prefix" : "https://staticsint.teams.cdn.office.net/evergreen-assets/safelinks/", "queryString": "url"}
		]
	}
}
```

### Rules

The top-level `rules` property allows you to define rules to configure a browser to be automatically launched for certain URLs.

The `browsers` property is an array of objects that match one of the following rule types.

Each rule configuration specifies a `browser` property that contains the `id` of the browser to launch if matched. This can also be `_prompt_` to force the list of browsers to be displayed.

#### URL Prefix Match

Performs a prefix match against the full URL.

| Name    | Type              | Description                                                              |
| ------- | ----------------- | ------------------------------------------------------------------------ |
| type    | string (required) | `prefix`                                                                 |
| prefix  | string (required) | The prefix to match                                                      |
| browser | string (required) | The `id` of the browser to launch or `_prompt_` to display the full list |

e.g.

```json
{
	"type": "prefix",
	"prefix": "https://dev.azure.com/myorg",
	"browser": "work"
}
```

#### Host Suffix Match

Perfoms a suffix match against the host portion of the URL. Handy for matching.

E.g. `www.github.com` and `github.com` would both match a rule of `github.com`.

| Name    | Type              | Description                                                              |
| ------- | ----------------- | ------------------------------------------------------------------------ |
| type    | string (required) | `host`                                                                   |
| host    | string (required) | The host suffix to match                                                 |
| browser | string (required) | The `id` of the browser to launch or `_prompt_` to display the full list |

e.g.

```json
{
	"type": "host",
	"prefix": "https://dev.azure.com/myorg",
	"browser": "work"
}
```

### Example configuration

```jsonc
{
	"browsers": [
		{
			"id": "iexplore",
			"name": "Internet Explorer",
			"exe": "C:\\Program Files\\Internet Explorer\\iexplore.exe",
			"iconPath": "C:\\Program Files\\Internet Explorer\\iexplore.exe",
			"hidden": true
		},
		{
			"id": "work",
			"name": "Microsoft Edge - Work",
			"exe": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"args": "--profile-directory=\u0022Default\u0022",
			"iconPath": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"hidden": false
		},
		{
			"id": "personal",
			"name": "Microsoft Edge - stuart@leeks.net (MSA)",
			"exe": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"args": "--profile-directory=\u0022Profile 1\u0022",
			"iconPath": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"hidden": false
		}
	],
	"rules": [
		{
			"type": "prefix",
			"prefix": "https://dev.azure.com/myorg",
			"browser": "work"
		},
		{
			"type": "host",
			"host": "github.com",
			"browser": "work"
		},
		{
			"type": "host",
			"host": "whatsapp.com",
			"browser": "personal"
		}
	]
}
```
