using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace pick_a_browser
{
    public static class Program
    {
        /// <summary>
        /// Application Entry Point.
        /// </summary>
        [STAThread]
        //public static async Task<int> Main(string[] args)
        public static int Main(string[] args)
        {
            App app = new App();
            app.InitializeComponent();

            if (args.Length > 0
                && args[0].Length >= 2
                && args[0].StartsWith("--"))
            {
                switch (args[0])
                {
                    case "--browser-scan":

                        var browsers = Browsers.Scan();
                        var window = new BrowserScanWindow(browsers);
                        return app.Run(window);
                }
            }

            return app.Run(new MainWindow());
        }
    }
}
