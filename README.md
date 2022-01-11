# pick-a-browser

`pick-a-browser` is a browser selector for Windows.
It registers as a browser and when launched it runs a configurable rule set to determine the real browser to launch.
If no rule matches then it prompts for the browser to launch from a configured list of browsers/profiles.

`pick-a-browser` was born after happily using (and contributing to) [BrowserPicker](https://github.com/mortenn/BrowserPicker).
After a number of reinstallations of the operating system, I wanted an automatic way to save my browser and rules configuration.
As a result, I created `pick-a-browser` as a way to address this (and to tweak behaviour in ways that suit my usage patterns slightly better).

## Installing

- Get the binaries
- Install `pick-a-browser` as a browser
- Create `pick-a-browser-settings.json`
- Set `pick-a-browser` as your default browser - https://support.microsoft.com/en-gb/windows/change-your-default-browser-in-windows-10-020c58c6-7d77-797a-b74e-8f07946c5db6

### Get the binaries

Currently, you can either grab the build artifact from [the latest CI build](https://github.com/stuartleeks/pick-a-browser/actions/workflows/ci-build.yml) and unzip, or clone and build from source.

### Installing pick-a-browser

To install `pick-a-browser`, run `pick-a-browser --install` (needs elevated permissions). You can use `pick-a-browser --uninstall` to uninstall

### Create settings file

By default, `pick-a-browser` will look for `pick-a-browser-settings.json` in your user profile folder and then in the same folder as the app itself.

If you which to put the settings in a different location, set the `PICK_A_BROWSER_CONFIG` environment variable to the full path to the settings file.

For details on configuring browsers/rules, see the [Configuration](#configuration) section. You can run `pick-a-browser --browser-scan` to generate the initial `browsers` section of the configuration.

## Configuration

TODO
- configuring browsers
- configuring rules
- specifying settings file location

