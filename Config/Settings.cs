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

namespace pick_a_browser.Config
{
	public class Settings
	{
		public Settings(List<Browser> browsers)
		{
			Browsers = browsers;
		}
		public List<Browser> Browsers { get; }

		// TODO - add tranformation rules (link shorteners, regex rules)
		// TODO - add browser matches (simple, prefix, regex)
	}

	public class Browser
	{
		public Browser(string name, string exe, string args, string iconPath)
		{
			Name = name;
			Exe = exe;
			Args = args;
			IconPath = iconPath;
		}
		public string Name { get; }
		public string Exe { get; }
		public string Args { get; }
		public string IconPath { get; }
	}

}
