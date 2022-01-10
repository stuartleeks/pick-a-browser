using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.IO;
using System.Net.NetworkInformation;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;

namespace pick_a_browser.Config
{
    public class Settings
    {
        public Settings(Browsers browsers, List<Rule> rules)
        {
            Browsers = browsers;
            Rules = rules;
        }
        public Browsers Browsers { get; }
        public List<Rule> Rules { get; }

        // TODO - add tranformation rules (link shorteners, regex rules)
        // TODO - add browser matches (simple, prefix, regex)
        public static async Task<Settings> LoadAsync(string filename)
        {
            return await SettingsSerialization.LoadFromFileAsync(GetSettingsFilename());
        }

        private static string GetSettingsFilename()
        {
            var settingsFilename = Environment.GetEnvironmentVariable("PICK_A_BROWSER_CONFIG");
            if (!string.IsNullOrEmpty(settingsFilename))
                return settingsFilename;

            var profilePath = Environment.GetEnvironmentVariable("USERPROFILE");
            if (string.IsNullOrEmpty(profilePath))
                throw new Exception("USERPROFILE not set");

            return Path.Join(profilePath, "pick-a-browser-settings.json");
        }

    }
}
