using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
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

namespace pick_a_browser
{
	/// <summary>
	/// Interaction logic for MainWindow.xaml
	/// </summary>
	public partial class MainWindow : Window
	{
		public MainWindow()
		{
			InitializeComponent();
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


    public class Config
    {
        
    }
}
