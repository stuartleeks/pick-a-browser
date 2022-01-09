using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System.Windows;

namespace pick_a_browser
{
	/// <summary>
	/// Interaction logic for App.xaml
	/// </summary>
	public partial class App : Application
	{

        public App()
        {
            Console.WriteLine("hi from app");
        }

		private string GetSettingsFilename()
		{
			var settingsFilename = Environment.GetEnvironmentVariable("PICK_A_BROWSER_CONFIG");
			if (!string.IsNullOrEmpty(settingsFilename))
				return settingsFilename;

			var profilePath = Environment.GetEnvironmentVariable("UESRPROFILE");
			if (string.IsNullOrEmpty(profilePath))
				throw new Exception("USERPROFILE not set");

			return Path.Join(profilePath, "pick-a-browser-settings.json");
		}
	}
}
