using Microsoft.Win32;
using System;
using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Security.Cryptography;
using System.Threading.Tasks;

namespace pick_a_browser.Config
{
    /// <summary>
    /// Represents a readonly list of Browser instances allowing for indexed access on the Id property
    /// </summary>
    public class Browsers : IEnumerable<Browser>, IReadOnlyList<Browser>
    {
        private readonly List<Browser> _browsers;

        public Browsers(List<Browser> browsers)
        {
            _browsers = browsers;
            // TODO - ensure that names are unique?
        }

        public Browser this[int index] => _browsers[index];
        public Browser this[string id] => _browsers.First(b => b.Id== id);

        public int Count => _browsers.Count;

        public IEnumerator<Browser> GetEnumerator()
        {
            return _browsers.GetEnumerator();
        }

        IEnumerator IEnumerable.GetEnumerator()
        {
            return GetEnumerator();
        }

        public static async Task<Browsers> Scan()
        {
            // Scan registered browsers as per: https://docs.microsoft.com/en-us/windows/win32/shell/start-menu-reg
            List<Browser> browserList = GetBrowsersFor(Registry.CurrentUser);
            browserList.AddRange(GetBrowsersFor(Registry.LocalMachine));

            var edge = browserList.FirstOrDefault(b => b.Name == "Microsoft Edge");
            if (edge != null)
            {
                browserList.Remove(edge);
                browserList.Add(new Browser(Guid.NewGuid().ToString(), edge.Name + " - Default", edge.Exe, "--profile-directory=\"Default\"", edge.IconPath, false));
                browserList.AddRange(
                    GetEdgeProfiles()
                    .Select(profile => new Browser(Guid.NewGuid().ToString(), edge.Name + " - " + profile, edge.Exe, $"--profile-directory=\"{profile}\"", edge.IconPath, false))
                );
            }

            try {
                // If we have existing settings, add new browsers at the end
                var existingSettings = await Settings.LoadAsync();
                var tmpList = existingSettings.Browsers.ToList();
                foreach (var browser in browserList)
                {
                    var existingBrowser = existingSettings.Browsers.FirstOrDefault(b=>b.Exe == browser.Exe && b.Args == browser.Args);
                    if (existingBrowser != null)
                        tmpList.Add(existingBrowser);
                }

                browserList = tmpList;
            }
            catch { }

            return new Browsers(browserList);
        }

        private static List<Browser> GetBrowsersFor(RegistryKey root)
        {
            var browsersKey = root.OpenSubKey("SOFTWARE\\Clients\\StartMenuInternet", writable: false);
            
            if (browsersKey == null)
                return new List<Browser>();

            return browsersKey.GetSubKeyNames()
                .Where(name => name != "pick-a-browser")
                .Select(name => BrowserFromRegistry(browsersKey.OpenSubKey(name)))
                .NonNulls()
                .ToList();
        }

        private static Browser? BrowserFromRegistry(RegistryKey? browserKey)
        {
            if (browserKey == null)
                return null;

            var exe = (string?)browserKey.OpenSubKey("shell\\open\\command", false)?.GetValue(null);
            if (exe == null)
                return null;

            if (exe.StartsWith('"') && exe.EndsWith('"'))
                exe = exe.Substring(1, exe.Length - 2);

            var name = (string?)browserKey.GetValue(null) ?? browserKey.Name;

            var iconPath = (string?)browserKey.OpenSubKey("DefaultIcon", false)?.GetValue(null);
            if (iconPath != null)
            {
                var iconCommaIndex = iconPath.IndexOf(',');
                if (iconCommaIndex > 0)
                {
                    iconPath = iconPath.Substring(0, iconCommaIndex);
                }
            }

            return new Browser(Guid.NewGuid().ToString(), name, exe, null, iconPath, false);
        }

        private static string[] GetEdgeProfiles()
        {
            var profilePath = Environment.GetEnvironmentVariable("USERPROFILE");

            var edgeUserDataPath = Path.Join(profilePath, "AppData\\Local\\Microsoft\\Edge\\User Data");

            var profilePrefix = Path.Join(edgeUserDataPath, "Profile ");
            var profileDirectories = Directory.GetDirectories(edgeUserDataPath)
                .Where(d => d.StartsWith(profilePrefix))
                .Select(d => d.Substring(edgeUserDataPath.Length + 1))
                .ToArray();
            return profileDirectories;
        }
    }
    public class Browser
    {
        public Browser(string id, string name, string exe, string? args, string? iconPath, bool hidden)
        {
            Id = id;
            Name = name;
            Exe = exe;
            Args = args;
            IconPath = iconPath;
            Hidden = hidden;
        }
        public string Id { get; }
        public string Name { get; }
        public string Exe { get; }
        public string? Args { get; }
        public string? IconPath { get; }
        public bool Hidden { get; set; }

        public void Launch(string url)
        {
            var args = Args == null ? url : $"{Args} {url}";
            Process.Start(Exe, args);
        }
    }

}
