using Microsoft.Win32;
using System;
using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;

namespace pick_a_browser.Config
{
    /// <summary>
    /// Represents a readonly list of Browser instances allowing for indexed access on the Name property
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
        public Browser this[string name] => _browsers.First(b => b.Name == name);

        public int Count => _browsers.Count;

        public IEnumerator<Browser> GetEnumerator()
        {
            return _browsers.GetEnumerator();
        }

        IEnumerator IEnumerable.GetEnumerator()
        {
            return GetEnumerator();
        }

        public static Browsers Scan()
        {
            // Scan registered browsers as per: https://docs.microsoft.com/en-us/windows/win32/shell/start-menu-reg
            var browsersKey = Registry.LocalMachine.OpenSubKey("SOFTWARE\\Clients\\StartMenuInternet", writable: false);

            if (browsersKey == null)
                return new Browsers(new List<Browser>());

            var browserList = browsersKey.GetSubKeyNames()
                .Where(name => name != "pick-a-browser")
                .Select(name => BrowserFromRegistry(browsersKey.OpenSubKey(name)))
                .NonNulls()
                .ToList();

            var edge = browserList.FirstOrDefault(b => b.Name == "Microsoft Edge");
            if (edge != null)
            {
                browserList.Remove(edge);
                browserList.Add(new Browser(edge.Name + " - Default", edge.Exe, "--profile-directory=\"Default\"", edge.IconPath));
                browserList.AddRange(
                    GetEdgeProfiles()
                    .Select(profile => new Browser(edge.Name + " - " + profile, edge.Exe, $"--profile-directory=\"{profile}\"", edge.IconPath))
                );
            }

            return new Browsers(browserList);
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

            return new Browser(name, exe, null, iconPath);
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
        public Browser(string name, string exe, string? args, string? iconPath)
        {
            Name = name;
            Exe = exe;
            Args = args;
            IconPath = iconPath;
        }
        public string Name { get; }
        public string Exe { get; }
        public string? Args { get; }
        public string? IconPath { get; }

        public void Launch(string url)
        {
            var args = Args == null ? url : $"{Args} {url}";
            Process.Start(Exe, args);
        }
    }

}
