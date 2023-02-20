# pick-a-browser

`pick-a-browser` is a browser selector for Windows.
It registers as a browser and when launched it runs a configurable rule set to determine the real browser to launch.
If no rule matches then it prompts for the browser to launch from a configured list of browsers/profiles.

`pick-a-browser` was born after happily using (and contributing to) [BrowserPicker](https://github.com/mortenn/BrowserPicker).
After a number of reinstallations of the operating system, I wanted an automatic way to save my browser and rules configuration.
As a result, I created `pick-a-browser` as a way to address this (and to tweak behaviour in ways that suit my usage patterns slightly better).

Contents:

- [Installing](#installing)
- [Configuration](#configuration)
- [Commands](#commands)

## Installing

-   Get the binaries
-   Install `pick-a-browser` as a browser
-   Create `pick-a-browser-settings.json`
-   Set `pick-a-browser` as your default browser - <https://support.microsoft.com/en-gb/windows/change-your-default-browser-in-windows-10-020c58c6-7d77-797a-b74e-8f07946c5db6>

### Get the binaries

Download `pick-a-browser.exe` from the [latest release](https://github.com/stuartleeks/pick-a-browser/releases/latest).

Alternatively,  you can either grab the build artifact from [the latest CI build](https://github.com/stuartleeks/pick-a-browser/actions/workflows/ci-build.yml) and unzip, or clone and build from source.

### Installing pick-a-browser

To install `pick-a-browser`, run `pick-a-browser --install` (needs elevated permissions). You can use `pick-a-browser --uninstall` to uninstall

### Create settings file

See [Configuration](#configuration) below.

## Configuration

By default, `pick-a-browser` will look for `pick-a-browser-settings.json` in your user profile folder and then in the same folder as the app itself.

If you wish to put the settings in a different location, set the `PICK_A_BROWSER_CONFIG` environment variable to the full path to the settings file.

### Updates

The `updates` property allows you to control the update behaviour of `pick-a-browser`.

The `updates` property can take any of the following values:

| Value              | Description                                                                                  |
|--------------------|----------------------------------------------------------------------------------------------|
| `none`             | `pick-a-browser` will not check for updates                                                  |
| `prompt` (default) | `pick-a-browser` will check for updates and display an indicator when an update is available |
| `auto`             | `pick-a-browser` will check for updates and automatically apply them in the background       |


### Logging

The `log` property allows you to control the logging behaviour of `pick-a-browser`.

Logs are written to `%LOCALAPPDATA%\StuartLeeks\pick-a-browser\pick-a-browser.log`

E.g. 

```json
{
	"log" : {
		"level": "info"
	}
}
```

#### Log level

The `level` property under `log` can take any of the following values:

| Value             | Description                                        |
|-------------------|----------------------------------------------------|
| `error` (default) | `pick-a-browser` will log errors                   |
| `info`            | `pick-a-browser` will log info messages and above  |
| `debug`           | `pick-a-browser` will log debug messages and above |

### Browsers

The top-level `browsers` property allows you to configure browsers (or browser profiles) that `pick-a-browser` should use.

You can run `pick-a-browser --browser-scan` to generate the initial `browsers` section of the configuration and then copy and paste this into the config file.

The `browsers` property is an array of objects with the following properties:

| Name       | Type              | Description                                                                  |
|------------|-------------------|------------------------------------------------------------------------------|
| `id`       | string (required) | The id to use to identify the browser - used to refer to browsers in rules   |
| `name`     | string (required) | The name to display in the browser picker UI                                 |
| `exe`      | string (required) | The path to the app to launch for the browser                                |
| `args`     | string (optional) | Any arguments to pass to the browser. Useful for specifying browser profiles |
| `iconPath` | string (optional) | not currently used                                                           |
| `hidden`   | bool (optional)   | When no rules are matched, the UI displays a list of all non-hidden browsers |

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

| Name          | Type              | Description                                              |
|---------------|-------------------|----------------------------------------------------------|
| `prefix`      | string (required) | The URL prefix to match for the shortener                |
| `queryString` | string (required) | The name of the query string value that contains the URL |

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

| Name      | Type              | Description                                                              |
|-----------|-------------------|--------------------------------------------------------------------------|
| `type`    | string (required) | `prefix`                                                                 |
| `prefix`  | string (required) | The prefix to match                                                      |
| `browser` | string (required) | The `id` of the browser to launch or `_prompt_` to display the full list |

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

| Name      | Type              | Description                                                              |
|-----------|-------------------|--------------------------------------------------------------------------|
| `type`    | string (required) | `host`                                                                   |
| `host`    | string (required) | The host suffix to match                                                 |
| `browser` | string (required) | The `id` of the browser to launch or `_prompt_` to display the full list |

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
	"updates" : "auto",
	"log" : {
		"level": "info"
	},
	"browsers": [
		{
			"id": "iexplore",
			"name": "Internet Explorer",
			"exe": "C:\\Program Files\\Internet Explorer\\iexplore.exe",
			"iconPath": "C:\\Program Files\\Internet Explorer\\iexplore.exe",
			"hidden": true // avoid being shown Internet Explorer as a choice even though it's registered in Windows
		},
		{
			"id": "work",
			"name": "Microsoft Edge - Work",
			"exe": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"args": "--profile-directory=\"Default\"",
			"iconPath": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"hidden": false
		},
		{
			"id": "personal",
			"name": "Microsoft Edge - Personal",
			"exe": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"args": "--profile-directory=\"Profile 1\"",
			"iconPath": "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe",
			"hidden": false
		}
	],
	"transformations" : {
		// NOTE that these examples are included in the default transformations
		"linkShorteners": [ "aka.ms" ], 
		"linkWrappers" : [
			{ "prefix" : "https://staticsint.teams.cdn.office.net/evergreen-assets/safelinks/", "queryString": "url"}
		]
	},
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

## Commands

This section documents the command line options for `pick-a-browser`:


### --install

`pick-a-browser --install` is used to install pick-a-browser, i.e. it registers `pick-a-browser` with Windows as a browser.

NOTE: This command needs to be run as administrator

### --uninstall

`pick-a-browser --uninstall` is used to uninstall pick-a-browser, i.e. unregister `pick-a-browser` with Windows.

NOTE: This command needs to be run as administrator

### --browser-scan

`pick-a-browser --browser-scan` can be used to generate the list of browsers to use in your configuration.

It lists the browsers registered with Windows and generates the configuration for `pick-a-browser`.

For Microsoft Edge and Google Chrome, this also includes any browser profiles that are configured.

This command can be re-run if you add a new browser or profile to re-generate the configuration.


### --update

`pick-a-browser --update` can be used to check for and apply any pending updates.

### &lt;url&gt;

`pick-a-browser <url>` is used to launch a browser with the specified url.

The configured transformations and rules will be applied to the url. If a rule matches at the end of this, that browser will be launched, otherwise the browser picker will be displayed.

NOTE: The url is optional and the selected browser will be launched with an empty page.

