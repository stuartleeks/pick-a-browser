using System;
using System.Buffers;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Net.Http.Json;
using System.Reflection;
using System.Text.Json.Serialization;
using System.Threading;
using System.Threading.Tasks;
using System.Windows;

namespace pick_a_browser
{
    public static class Updater
    {
        public static IDisposable GetUpdateLock()
        {
            var filename = GetUpdateLockFilename();
            var path = Path.GetDirectoryName(filename);
            if (path != null)
                Directory.CreateDirectory(path);

            try
            {
                var handle = File.OpenHandle(filename, FileMode.Create, FileAccess.Write, FileShare.None);
                return handle;
            }
            catch (IOException ioe) when (ioe.HResult == -2147024864) // SHARING_VIOLATION - 0x80070020
            {
                throw new UnableToObtainUpdateLockException("Failed to obtain update lock", ioe);
            }
        }
        private static string GetUpdateLockFilename()
        {
            return Path.Join(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
                            "StuartLeeks\\pick-a-browser\\update-lock.txt");
        }

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

        public static async Task<Version?> CheckForUpdateAsync(CancellationToken cancellationToken)
        {
            var githubRelease = await GetLatestGitHubReleaseAsync(cancellationToken);
            if (githubRelease == null)
            {
                MessageBox.Show("Failed to get latest GitHub release");
                return null;
            }
            if (githubRelease.TagName == null)
            {
                MessageBox.Show("Failed to get latest GitHub version");
                return null;
            }
            var versionString = githubRelease.TagName.Trim('v', 'V');
            var githubVersion = new Version(versionString);

            var currentVersion = GetAssemblyVersion();

            var updateAvailable = githubVersion.CompareTo(currentVersion) > 0;

            if (updateAvailable)
                return githubVersion;
            return null;
        }

        public static async Task UpdateAsync(CancellationToken cancellationToken, Action<string>? statusUpdater)
        {
            using var _ = GetUpdateLock(); // Ensure we only have a single updater running

            statusUpdater?.Invoke("Getting release details...");

            var githubRelease = await GetLatestGitHubReleaseAsync(cancellationToken);
            if (githubRelease == null)
                throw new Exception("Failed to get latest GitHub release");
            if (githubRelease.TagName == null)
                throw new Exception("Failed to get latest GitHub version");

            var versionString = githubRelease.TagName.Trim('v', 'V');
            var githubVersion = new Version(versionString);

            var currentVersion = GetAssemblyVersion();

            var updateAvailable = githubVersion.CompareTo(currentVersion) > 0;

            if (!updateAvailable)
                return;

            var asset = githubRelease.Assets.FirstOrDefault(a => a.Name == "pick-a-browser.exe");
            if (asset == null)
                throw new Exception("Release didn't contain pick-a-browser.exe asset");

            var tmpExePath = Path.Join(AppContext.BaseDirectory, "tmp.pick-a-browser.exe");  // Can't use Assembly.GetExecutingAssembly().Location in Single File App

            cancellationToken.ThrowIfCancellationRequested();

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

            statusUpdater?.Invoke("Downloading new release...");
            using (var stream = await client.GetStreamAsync(asset.DownloadUrl, cancellationToken))
            using (var fileStream = File.Create(tmpExePath))
            {
                await CopyToWithProgress(stream, fileStream, statusUpdater, cancellationToken);
            }

            cancellationToken.ThrowIfCancellationRequested();

            var exePath = Path.Join(AppContext.BaseDirectory, "pick-a-browser.exe");
            var oldExePath = Path.Join(AppContext.BaseDirectory, "old.pick-a-browser.exe");

            statusUpdater?.Invoke($"Downloaded complete");
            statusUpdater?.Invoke("Renaming old exe...");
            if (File.Exists(oldExePath))
                File.Delete(oldExePath);
            File.Move(exePath, oldExePath);
            statusUpdater?.Invoke("Replacing with downloaded exe...");
            File.Move(tmpExePath, exePath);

            statusUpdater?.Invoke("Done!");
        }

        private static async Task CopyToWithProgress(Stream source, Stream destination, Action<string>? statusUpdater, CancellationToken cancellationToken)
        {
            long totalBytesRead = 0;
            long lastReported = 0;
            long reportIncrements = 10 * 1024 * 1024; // 10 MB

            // From https://github.com/dotnet/coreclr/blob/a6045809b3720833459c9247aeb4aafe281d7841/src/mscorlib/src/System/IO/Stream.cs#L125-L152
            //We pick a value that is the largest multiple of 4096 that is still smaller than the large object heap threshold (85K).
            // The CopyTo/CopyToAsync buffer is short-lived and is likely to be collected at Gen0, and it offers a significant
            // improvement in Copy performance.
            const int CopyBufferSize = 81920;
            int bufferSize = CopyBufferSize;
            byte[] buffer = ArrayPool<byte>.Shared.Rent(bufferSize);
            bufferSize = 0; // reuse same field for high water mark to avoid needing another field in the state machine
            try
            {
                while (true)
                {
                    int bytesRead = await source.ReadAsync(buffer, 0, buffer.Length, cancellationToken).ConfigureAwait(false);
                    if (bytesRead == 0) break;
                    if (bytesRead > bufferSize) bufferSize = bytesRead;
                    await destination.WriteAsync(buffer, 0, bytesRead, cancellationToken).ConfigureAwait(false);
                    totalBytesRead += bytesRead;
                    if (totalBytesRead - lastReported > reportIncrements)
                    {
                        statusUpdater?.Invoke($"Downloaded {totalBytesRead / (1024 * 1024)} MB...");
                        lastReported = totalBytesRead;
                    }
                }
            }
            finally
            {
                Array.Clear(buffer, 0, bufferSize); // clear only the most we used
                ArrayPool<byte>.Shared.Return(buffer, clearArray: false);
            }
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
        public static async Task<GitHubRelease?> GetLatestGitHubReleaseAsync(CancellationToken cancellationToken)
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
            var httpResponse = await client.GetAsync("https://api.github.com/repos/stuartleeks/pick-a-browser/releases/latest", cancellationToken);

            var response = await httpResponse.Content.ReadFromJsonAsync<GitHubRelease>(cancellationToken: cancellationToken);

            return response;
        }
    }


    [Serializable]
    public class UnableToObtainUpdateLockException : Exception
    {
        public UnableToObtainUpdateLockException() { }
        public UnableToObtainUpdateLockException(string message) : base(message) { }
        public UnableToObtainUpdateLockException(string message, Exception inner) : base(message, inner) { }
        protected UnableToObtainUpdateLockException(
          System.Runtime.Serialization.SerializationInfo info,
          System.Runtime.Serialization.StreamingContext context) : base(info, context) { }
    }
}
