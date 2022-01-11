using pick_a_browser.Config;
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
	/// Interaction logic for LoadingWindow.xaml
	/// </summary>
	public partial class LoadingWindow: Window
	{
		public LoadingWindow(LoadingViewModel viewModel)
		{
			InitializeComponent();
            DataContext = viewModel;
		}
	}

    public class LoadingViewModel : ViewModel
    {
        private string _url="";
        public string Url
        {
            get { return _url; }
            set { _url = value; FirePropertyChanged(); }
        }
    }

    public class DesignTimeLoadingViewModel: LoadingViewModel
    {
        public DesignTimeLoadingViewModel()
        {
            Url = "https://example.com/some/path/goes/here";
        }
    }
}
