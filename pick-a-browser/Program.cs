using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace pick_a_browser
{
    public static class Program
    {
        /// <summary>
        /// Application Entry Point.
        /// </summary>
        public static async Task Main(string[] args)
        {
            if (args.Length > 0
                && args[0].Length >= 2
                && args[0].StartsWith("--"))
            {
                switch (args[0])
                {
                    case "--browser-scan":
                        RunWpfApp(() =>
                        {
                            App app = new App();
                            app.InitializeComponent();
                            var scannedBrowsers = Browsers.Scan();
                            var window = new BrowserScanWindow(scannedBrowsers);
                            app.Run(window);
                        });
                        return;
                }
            }

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

            RunWpfApp(() =>
            {
                App app = new App();
                app.InitializeComponent();
                var model = new PickABrowserViewModel(browsers, url);
                var window = new PickABrowserWindow(model);
                app.Run(window);
            });
        }

        private static void RunWpfApp(ThreadStart action)
        {
            // https://github.com/dotnet/roslyn/issues/37122
            var thread = new Thread(action);
            thread.SetApartmentState(ApartmentState.STA);
            thread.Start();
            thread.Join();
            return;
        }
    }
}
