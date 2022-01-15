using Microsoft.Win32;
using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.IO;
using System.Linq;
using System.Net.Http;
using System.Reflection;
using System.Threading;
using System.Threading.Tasks;
using System.Web;
using System.Windows;

namespace pick_a_browser
{
    /// <summary>
    /// Interaction logic for App.xaml
    /// </summary>
    public partial class App : Application
    {
        protected async override void OnStartup(StartupEventArgs e)
        {
            try
            {
                var args = e.Args;

                if (args.Length > 0
                 && args[0].Length >= 2
                 && args[0].StartsWith("--"))
                {
                    switch (args[0])
                    {
                        case "--browser-scan":
                            await RunBrowserScanAsync();
                            return;

                        case "--install":
                            RunInstall();
                            Current.Shutdown();
                            return;

                        case "--uninstall":
                            RunUninstall();
                            Current.Shutdown();
                            return;

                        case "--update":
                            await RunUpdateAsync();
                            Current.Shutdown();
                            return;
                    }
                }

                await RunPickABrowserAsync(args);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
                Current.Shutdown();
            }
        }

        private static async Task RunBrowserScanAsync()
        {
            var scannedBrowsers = await Browsers.Scan();
            var window = new BrowserScanWindow(scannedBrowsers);
            window.Show();
        }

        private static void RunInstall()
        {
            // Register as browser as per: https://docs.microsoft.com/en-us/windows/win32/shell/start-menu-reg
            var root = Registry.LocalMachine;
            var browsersKey = root.CreateSubKey("SOFTWARE\\Clients\\StartMenuInternet", writable: true);
            if (browsersKey == null)
            {
                MessageBox.Show("Failed to open SOFTWARE\\Clients\\StartMenuInternet");
                return;
            }

            var pickABrowserKey = browsersKey.OpenSubKey("pick-a-browser");
            if (pickABrowserKey != null)
            {
                MessageBox.Show("pick-a-browser already installed");
                return;
            }

            // TODO - look at handling renamed exe file
            var exePath = Path.Join(AppContext.BaseDirectory, "pick-a-browser.exe");  // Can't use Assembly.GetExecutingAssembly().Location in Single File App

            pickABrowserKey = browsersKey.CreateSubKey("pick-a-browser", writable: true);
            pickABrowserKey.SetValue(null, "Pick A Browser");

            pickABrowserKey.CreateSubKey("DefaultIcon", writable: true).SetValue(null, $"{exePath},0");

            var capabilitiesKey = pickABrowserKey.CreateSubKey("Capabilities", writable: true);
            capabilitiesKey.SetValue("ApplicationDescription", "browser selector - see https://github.com/stuartleeks/pick-a-browser");
            capabilitiesKey.SetValue("ApplicationName", "pick-a-browser");
            capabilitiesKey.CreateSubKey("StartMenu").SetValue("StartMenuInternet", "pick-a-browser");
            capabilitiesKey.CreateSubKey("FileAssociations", writable: true);
            var urlAssociationsKey = capabilitiesKey.CreateSubKey("UrlAssociations", writable: true);
            urlAssociationsKey.SetValue("http", "pick-a-browser");
            urlAssociationsKey.SetValue("https", "pick-a-browser");


            // InstallInfo: https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/cc144109(v=vs.85)
            var installInfoKey = pickABrowserKey.CreateSubKey("InstallInfo", writable: true);
            installInfoKey.SetValue("HideIconsCommand", "");
            installInfoKey.SetValue("ReinstallCommand", "");
            installInfoKey.SetValue("ShowIconsCommand", "");
            installInfoKey.SetValue("IconsVisible", "1");

            var shellKey = pickABrowserKey.CreateSubKey("shell", writable: true);
            shellKey.SetValue(null, "open");
            var commandKey = shellKey.CreateSubKey("open\\command", writable: true);
            commandKey.SetValue(null, exePath); // no param here

            // https://docs.microsoft.com/en-us/windows/win32/shell/default-programs#registering-an-application-for-use-with-default-programs
            // HKCR:
            var classesPickABrowserKey = Registry.ClassesRoot.CreateSubKey("pick-a-browser", writable: true);
            classesPickABrowserKey.SetValue(null, "pick-a-browser");
            classesPickABrowserKey.CreateSubKey("DefaultIcon", writable: true).SetValue(null, $"{exePath},0");

            shellKey = classesPickABrowserKey.CreateSubKey("shell", writable: true);
            shellKey.SetValue(null, "open");
            commandKey = shellKey.CreateSubKey("open\\command", writable: true);
            commandKey.SetValue(null, $"\"{exePath}\" \"%1\"");

            var classesAppKey = classesPickABrowserKey.CreateSubKey("Application", writable: true);
            classesAppKey.SetValue("ApplicationDescription", "browser selector - see https://github.com/stuartleeks/pick-a-browser");
            classesAppKey.SetValue("ApplicationName", "pick-a-browser");
            classesAppKey.SetValue("ApplicationIcon", $"{exePath},0");
            classesPickABrowserKey.CreateSubKey("DefaultIcon", writable: true).SetValue(null, $"{exePath},0");

            capabilitiesKey = classesPickABrowserKey.CreateSubKey("Capabilities", writable: true);
            capabilitiesKey.SetValue("ApplicationDescription", "browser selector - see https://github.com/stuartleeks/pick-a-browser");
            capabilitiesKey.SetValue("ApplicationName", "pick-a-browser");
            capabilitiesKey.CreateSubKey("Startmenu").SetValue("StartmenuInternet", "pick-a-browser");
            capabilitiesKey.CreateSubKey("FileAssociations", writable: true);
            urlAssociationsKey = capabilitiesKey.CreateSubKey("UrlAssociations", writable: true);
            urlAssociationsKey.SetValue("http", "pick-a-browser");
            urlAssociationsKey.SetValue("https", "pick-a-browser");


            // Software\p-a-b
            capabilitiesKey = root.CreateSubKey("SOFTWARE\\pick-a-browser\\Capabilities", writable: true);
            capabilitiesKey.SetValue("ApplicationDescription", "browser selector - see https://github.com/stuartleeks/pick-a-browser");
            capabilitiesKey.SetValue("ApplicationName", "pick-a-browser");
            capabilitiesKey.CreateSubKey("StartMenu").SetValue("StartMenuInternet", "pick-a-browser");
            capabilitiesKey.CreateSubKey("FileAssociations", writable: true);
            urlAssociationsKey = capabilitiesKey.CreateSubKey("UrlAssociations", writable: true);
            urlAssociationsKey.SetValue("http", "pick-a-browser");
            urlAssociationsKey.SetValue("https", "pick-a-browser");

            root.CreateSubKey("SOFTWARE\\RegisteredApplications", writable: true).SetValue("pick-a-browser", "SOFTWARE\\Clients\\StartmenuInternet\\pick-a-browser\\Capabilities");

            MessageBox.Show("pick-a-browser installed");
        }

        private static void RunUninstall()
        {
            var root = Registry.LocalMachine;
            var browsersKey = root.OpenSubKey("SOFTWARE\\Clients\\StartMenuInternet", writable: true);
            if (browsersKey == null)
            {
                MessageBox.Show("Failed to open SOFTWARE\\Clients\\StartMenuInternet");
                return;
            }

            root.OpenSubKey("SOFTWARE", writable: true)?.DeleteSubKeyTree("pick-a-browser", throwOnMissingSubKey: false);
            Registry.ClassesRoot.DeleteSubKeyTree("pick-a-browser", throwOnMissingSubKey: false);
            root.CreateSubKey("SOFTWARE\\RegisteredApplications", writable: true).DeleteValue("pick-a-browser", throwOnMissingValue: false);

            var pickABrowserKey = browsersKey.OpenSubKey("pick-a-browser", writable: true);
            if (pickABrowserKey == null)
            {
                MessageBox.Show("pick-a-browser not installed");
                return;
            }

            browsersKey.DeleteSubKeyTree("pick-a-browser");
            MessageBox.Show("pick-a-browser uninstalled");
        }
        private static async Task RunUpdateAsync()
        {
            try
            {
                await Updater.UpdateAsync(null);
                MessageBox.Show("Done!");
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Error: {ex.Message}");
            }
        }

        private static async Task RunPickABrowserAsync(string[] args)
        {
            var appData = await AppData.Load();
            var settings = await Settings.LoadAsync();

            var originalUrl = args.Length > 0 ? args[0] : "";
            var loadingViewModel = new LoadingViewModel { Url = originalUrl };

            var cts = new CancellationTokenSource();
            var loadingTask = ShowLoadingAsync(loadingViewModel, cts.Token);
            var url = await UnwrapLinkAsync(originalUrl, u => loadingViewModel.Url = u, settings.Transformations, cts.Token);
            cts.Cancel();

            // Get matches with highest weights (handle multiple matches with the same weight)
            var uri = new Uri(url);
            var matchedBrowserIds = settings.Rules
                    .Select(r => r.GetMatch(uri))
                    .Where(m => m.Weight > 0)
                    .GroupBy(m => m.Weight)
                    .OrderByDescending(g => g.Key)
                    // Take the top weighted match(es)
                    ?.FirstOrDefault()
                    // Get browsers
                    ?.Select(m => m.BrowserId)
                    ?.Distinct()
                    .ToList();

            List<Browser>? browsers;
            if (matchedBrowserIds == null || matchedBrowserIds.Count == 0 || matchedBrowserIds[0] == "_prompt_")
            {
                browsers = settings.Browsers.Where(b => !b.Hidden).ToList();
            }
            else
            {
                browsers = settings.Browsers.Where(b => matchedBrowserIds.Contains(b.Id)).ToList();
            }

            if (browsers.Count == 1)
            {
                browsers[0].Launch(url);
                Current.Shutdown();
            }


            var model = new PickABrowserViewModel(browsers, originalUrl, url);
            if (DateTime.UtcNow - appData.LastUpdateCheckUtc > TimeSpan.FromHours(4)) // TODO - config for frequency?
            {
                // Check GitHub
                // IIFE style approach to async execution without awaiting
                var _ = ((Func<Task>)(async () =>
                {
                    appData.LastUpdateCheckUtc = DateTime.UtcNow;
                    await appData.SaveAsync();
                    appData.LastUpdateCheckGitHubVersion = await Updater.CheckForUpdateAsync();
                    await appData.SaveAsync();
                    model.UpdateAvailable = appData.LastUpdateCheckGitHubVersion;
                }))();
            }
            else
            {
                // Use cached version from last check
                var currentVersion = Updater.GetAssemblyVersion();
                if ((appData.LastUpdateCheckGitHubVersion?.CompareTo(currentVersion) ?? -1) > 0)
                    model.UpdateAvailable = appData.LastUpdateCheckGitHubVersion;
            }
            var window = new PickABrowserWindow(model);
            window.Show();
        }

        private static async Task ShowLoadingAsync(LoadingViewModel viewModel, CancellationToken cancellationToken)
        {
            await Task.Delay(300); // TODO add to config?

            if (cancellationToken.IsCancellationRequested)
                return;

            var window = new LoadingWindow(viewModel);

            // Hide when cancelled
            cancellationToken.Register(() => window.Hide());

            window.Show(); // TODO - allow cancelling from the LoadingWindow?
        }

        private static async Task<string> UnwrapLinkAsync(string url, Action<string> urlUpdated, Transformations transformations, CancellationToken cancellationToken)
        {
            var result = url;
            var client = new HttpClient(new HttpClientHandler { AllowAutoRedirect = false });

            var linkShorteners = DefaultLinkShorteners.Concat(transformations.LinkShorteners).ToList();
            var linkWrappers = DefaultLinkWrappers.Concat(transformations.LinkWrappers).ToList();

            while (true)
            {
                var uri = new Uri(result);
                var shortener = linkShorteners.FirstOrDefault(s => uri.Host == s);
                if (shortener != null)
                {
                    var response = await client.GetAsync(uri, cancellationToken);
                    var newUrl = response.Headers.Location?.OriginalString;
                    if (newUrl == null)
                    {
                        // Can't process any further
                        return result;
                    }
                    result = newUrl;
                    continue;
                }

                var wrapper = linkWrappers.FirstOrDefault(s => result.StartsWith(uri.AbsoluteUri));
                if (wrapper != null)
                {
                    var query = HttpUtility.ParseQueryString(uri.Query);
                    var queryValue = query[wrapper.QueryStringKey];
                    if (queryValue == null)
                    {
                        // Can't process any further
                        return result;
                    }
                    result = queryValue;
                    continue;
                }

                break; // didn't match shortener or wrappers
            }
            return result;
        }

        /// <summary>
        /// linkShorteners require a GET request which returns a redirect resopnse with Location header for the wrapped URL
        /// </summary>
        private static readonly string[] DefaultLinkShorteners = // TODO - allow adding via config
        {
            "aka.ms",
            "t.co",
            "go.microsoft.com",
        };

        /// <summary>
        /// linkWrappers contain the target URL in a query string value
        /// </summary>
        private static readonly LinkWrapper[] DefaultLinkWrappers = {
            new LinkWrapper("https://staticsint.teams.cdn.office.net/evergreen-assets/safelinks/", "url"),
            new LinkWrapper("https://nam06.safelinks.protection.outlook.com/", "url"),
        };

    }
}
