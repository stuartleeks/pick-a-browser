using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;

namespace pick_a_browser
{
    /// <summary>
    /// Interaction logic for BrowserScanWindow.xaml
    /// </summary>
    public partial class BrowserScanWindow : Window
    {
        public BrowserScanWindow(Browsers browsers)
        {
            InitializeComponent();

            var browserJson = SettingsSerialization.ToJsonNode(browsers)
                                .ToJsonString(
                                    new JsonSerializerOptions
                                    {
                                        WriteIndented = true
                                    });

            browserContent.Text = browserJson;
        }

    }

}
