using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Net.Http.Json;
using System.Reflection;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using System.Windows;

namespace pick_a_browser
{
    public class Updater
    {

        public static Version? GetAssemblyVersion()
        {
            var assembly = typeof(pick_a_browser.App)!.Assembly;
            return assembly.GetName()?.Version;
        }

        public static AssemblyInformationalVersionAttribute? GetAssemblyInformationalVersion()
        {
            var assembly = typeof(pick_a_browser.App)!.Assembly;
            return Attribute.GetCustomAttribute(assembly, typeof(AssemblyInformationalVersionAttribute)) as AssemblyInformationalVersionAttribute;
        }

        public static async Task UpdateAsync()
        {
            var githubRelease = await Updater.GetLatestGitHubReleaseAsync();
            if (githubRelease == null)
            {
                MessageBox.Show("Failed to get latest GitHub release");
                return;
            }
            if (githubRelease.TagName == null)
            {
                MessageBox.Show("Failed to get latest GitHub version");
                return;
            }
            var versionString = githubRelease.TagName.Trim('v', 'V');
            var githubVersion = new Version(versionString);

            var currentVersion = GetAssemblyVersion();

            var updateAvailable = githubVersion.CompareTo(currentVersion) > 0;

            if (!updateAvailable)
                return;

            var asset = githubRelease.Assets.FirstOrDefault(a => a.Name == "pick-a-browser.exe");
            if (asset == null)
            {
                MessageBox.Show("Release didn't contain pick-a-browser.exe asset");
                return;
            }

            var tmpExePath = Path.Join(AppContext.BaseDirectory, "tmp.pick-a-browser.exe");  // Can't use Assembly.GetExecutingAssembly().Location in Single File App

            var client = new HttpClient
            {
                DefaultRequestHeaders =
                {
                    UserAgent =
                    {
                        new ProductInfoHeaderValue("pick-a-browser", GetAssemblyVersion()?.ToString())
                    }
                }
            };

            using (var stream = await client.GetStreamAsync(asset.DownloadUrl))
            using (var fileStream = File.OpenWrite(tmpExePath))
            {
                // TODO - show update progress
                await stream.CopyToAsync(fileStream);
            }


            var exePath = Path.Join(AppContext.BaseDirectory, "pick-a-browser.exe");
            var oldExePath = Path.Join(AppContext.BaseDirectory, "old.pick-a-browser.exe");

            File.Move(exePath, oldExePath);
            File.Move(tmpExePath, exePath);

            MessageBox.Show("Done!");
        }

        public class GitHubRelease
        {
            [JsonPropertyName("tag_name")]
            public string? TagName { get; set; }

            [JsonPropertyName("assets")]
            public List<GitHubReleaseAsset> Assets { get; set; } = new List<GitHubReleaseAsset>();
        }
        public class GitHubReleaseAsset
        {
            [JsonPropertyName("browser_download_url")]
            public string? DownloadUrl { get; set; }

            [JsonPropertyName("name")]
            public string? Name { get; set; }
        }
        public static async Task<GitHubRelease?> GetLatestGitHubReleaseAsync()
        {
            var client = new HttpClient
            {
                DefaultRequestHeaders =
                {
                    Accept =
                    {
                        new MediaTypeWithQualityHeaderValue("application/vnd.github.v3+json")
                    },
                    UserAgent =
                    {
                        new ProductInfoHeaderValue("pick-a-browser", GetAssemblyVersion()?.ToString())
                    }
                }
            };

            // From https://docs.github.com/en/rest/reference/releases#get-the-latest-release
            //      The latest release is the most recent non-prerelease, non-draft release,
            //      sorted by the created_at attribute.
            //      The created_at attribute is the date of the commit used for the release,
            //      and not the date when the release was drafted or published.
            var httpResponse = await client.GetAsync("https://api.github.com/repos/stuartleeks/pick-a-browser/releases/latest");

            var response = await httpResponse.Content.ReadFromJsonAsync<GitHubRelease>();

            return response;
        }
    }
}
