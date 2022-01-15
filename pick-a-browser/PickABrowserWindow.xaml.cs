using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Reflection;
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
                Close();
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
            _browsers = browsers.Select((b, i) => new BrowserViewModel(b, url, i)).ToList();
            _originalUrl = originalUrl;
            _url = url;

            Version = Updater.GetAssemblyVersion()?.ToString();
            InformationalVersion = Updater.GetAssemblyInformationalVersion()?.InformationalVersion;
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

        private Version? _updateAvailable;
        public Version? UpdateAvailable
        {
            get { return _updateAvailable; }
            set
            {
                _updateAvailable = value;
                FirePropertyChanged();
                FirePropertyChanged(nameof(UpdateAvailableVisibility));
            }
        }
        public Visibility UpdateAvailableVisibility
        {
            // TODO - use converter from UpdateAvailable property
            get => _updateAvailable == null ? Visibility.Collapsed : Visibility.Visible;
        }

        public string? Version { get; }
        public string? InformationalVersion { get; }



        private bool _isUpdating = false;
        public DelegateCommand<object?> Update => new DelegateCommand<object?>(foo =>
        {
            _isUpdating = true;
            var viewModel = new UpdateViewModel(UpdateAvailable!, autoStart: false);
            var window = new UpdateWindow(viewModel);
            window.Show();
        }, _ => UpdateAvailable != null && !_isUpdating);
    }

    public class DesignTimePickABrowserViewModel : PickABrowserViewModel
    {
        public DesignTimePickABrowserViewModel()
            : base(GetBrowsers(), "https://aka.ms/example", "https://example.com/some/path")
        {
            UpdateAvailable = new Version("1.2.3");
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
        private readonly int _index;

        public BrowserViewModel(Browser browser, string url, int index)
        {
            _browser = browser;
            _url = url;
            _index = index;
        }

        public string Name { get => _browser.Name; }
        public string? IconPath { get => _browser.IconPath; }

        public string DisplayText { get => $"{_index + 1}: {Name}"; }
        public DelegateCommand<Window?> Launch => new DelegateCommand<Window?>(window =>
        {
            _browser.Launch(_url);
            if (window != null)
                window.Close();
        });
    }
}
