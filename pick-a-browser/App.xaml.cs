using Microsoft.Win32;
using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;
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
                            RunBrowserScan();
                            return;

                        case "--install":
                            RunInstall();
                            Current.Shutdown();
                            return;

                        case "--uninstall":
                            RunUninstall();
                            Current.Shutdown();
                            return;
                    }
                }

                await RunPickABrowser(args);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
                Current.Shutdown();
            }
        }

        private static void RunBrowserScan()
        {
            var scannedBrowsers = Browsers.Scan();
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

        private static async Task RunPickABrowser(string[] args)
        {
            var settings = await Settings.LoadAsync();

            var browsers = settings.Browsers.ToList();
            var url = args.Length > 0 ? args[0] : "";

            if (url != "")
            {
                var uri = new Uri(url);

                // Get matches with highest weights (handle multiple matches with the same weight)
                var matchedBrowserNames = settings.Rules
                        .Select(r => r.GetMatch(uri))
                        .Where(m => m.Weight > 0)
                        .GroupBy(m => m.Weight)
                        .OrderByDescending(g => g.Key)
                        // Take the top weighted match(es)
                        ?.FirstOrDefault()
                        // Get browsers
                        ?.Select(m => m.BrowserName)
                        ?.Distinct();

                if (matchedBrowserNames != null)
                {
                    browsers = browsers.Where(b => matchedBrowserNames.Contains(b.Name)).ToList();
                }
            }

            if (browsers.Count == 1)
            {
                browsers[0].Launch(url);
                Current.Shutdown();
            }

            var model = new PickABrowserViewModel(browsers, url);
            var window = new PickABrowserWindow(model);
            window.Show();
        }

    }
}
