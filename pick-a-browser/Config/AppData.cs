using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

namespace pick_a_browser.Config
{
    public class AppData
    {
        private static AppData GetDefault()
        {
            return new AppData
            {
                LastUpdateCheckUtc = DateTime.MinValue
            };
        }
        public DateTime LastUpdateCheckUtc { get; set; }
        public Version? LastUpdateCheckGitHubVersion { get; set; }

        public async Task SaveAsync()
        {
            await SaveToFileAsync(GetAppDataFilename(), this);
        }

        public static async Task<AppData> Load()
        {
            return await LoadFromFileAsync(GetAppDataFilename()) ?? GetDefault();
        }

        private static string GetAppDataFilename()
        {
            return Path.Join(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
                            "StuartLeeks\\pick-a-browser\\pick-a-browser-data.json");
        }

        public static async Task<AppData?> LoadFromFileAsync(string filename)
        {
            if (!File.Exists(filename))
                return null;

            using (var stream = File.OpenRead(filename))
                return await JsonSerializer.DeserializeAsync<AppData>(stream);
        }

        public static async Task SaveToFileAsync(string filename, AppData appData)
        {
            var path = Path.GetDirectoryName(filename);
            if (path != null)
                Directory.CreateDirectory(path);


            using (var stream = File.Create(filename))
                await JsonSerializer.SerializeAsync(stream, appData);
        }
    }
}
