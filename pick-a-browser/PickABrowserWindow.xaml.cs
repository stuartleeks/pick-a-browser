using pick_a_browser.Config;
using System.Collections.Generic;
using System.Diagnostics;
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
    /// Interaction logic for PickABrowserWindow.xaml
    /// </summary>
    public partial class PickABrowserWindow : Window
    {
        private PickABrowserViewModel _viewModel;
        public PickABrowserWindow(PickABrowserViewModel viewModel)
        {
            InitializeComponent();
            _viewModel = viewModel;
            DataContext = _viewModel;
        }


        private void Window_KeyDown(object sender, KeyEventArgs e)
        {
            switch (e.Key)
            {
                case Key.C:
                    Clipboard.SetText(_viewModel.Url);
                    return;
                case Key.Escape:
                    Close();
                    return;
            }
            var browserIndex = GetIndexFromKey(e.Key);
            if (browserIndex != null && browserIndex < _viewModel.Browsers.Count)
            {
                _viewModel.Browsers[(int)browserIndex].Launch.Execute(null);
                Application.Current.Shutdown();
                return;
            }
        }
        private int? GetIndexFromKey(Key key)
        {
            switch (key)
            {
                case Key.D1: return 0;
                case Key.D2: return 1;
                case Key.D3: return 2;
                case Key.D4: return 3;
                case Key.D5: return 4;
                case Key.D6: return 5;
                case Key.D7: return 6;
                case Key.D8: return 7;
                case Key.D9: return 8;
                case Key.D0: return 9;
                default: return null;
            }
        }
    }
    public class PickABrowserViewModel : ViewModel
    {
        public PickABrowserViewModel(List<Browser> browsers, string originalUrl, string url)
        {
            _browsers = browsers.Select(b => new BrowserViewModel(b, url)).ToList();
            _originalUrl = originalUrl;
            _url = url;
        }

        private List<BrowserViewModel> _browsers;
        public List<BrowserViewModel> Browsers
        {
            get { return _browsers; }
            set { _browsers = value; FirePropertyChanged(); }
        }


        private string _originalUrl;
        public string OriginalUrl
        {
            get { return _originalUrl; }
            set { _originalUrl = value; FirePropertyChanged(); }
        }

        private string _url;
        public string Url
        {
            get { return _url; }
            set { _url = value; FirePropertyChanged(); }
        }
    }

    public class DesignTimePickABrowserViewModel : PickABrowserViewModel
    {
        public DesignTimePickABrowserViewModel()
            : base(GetBrowsers(), "https://aka.ms/example", "https://example.com/some/path")
        {
        }

        private static List<Browser> GetBrowsers()
        {
            return new List<Browser>
            {
                new Browser("test1", "Browser number one", "", null, null, false),
                new Browser("test2", "Browser number two", "", null, null, false),
                new Browser("test3", "Browser number three", "", null, null, false),
                new Browser("test4", "Browser number four", "", null, null, false),
            };
        }
    }

    public class BrowserViewModel
    {
        private readonly Browser _browser;
        private readonly string _url;

        public BrowserViewModel(Browser browser, string url)
        {
            _browser = browser;
            _url = url;
        }

        public string Name { get => _browser.Name; }
        public string? IconPath { get => _browser.IconPath; }

        public DelegateCommand<object?> Launch => new DelegateCommand<object?>(foo =>
        {
            _browser.Launch(_url);
            Application.Current.Shutdown();
        });
    }
}
