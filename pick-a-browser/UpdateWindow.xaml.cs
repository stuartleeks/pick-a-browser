using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Threading;
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
    /// Interaction logic for UpdateWindow.xaml
    /// </summary>
    public partial class UpdateWindow : Window
    {
        private readonly UpdateViewModel _viewModel;

        public UpdateWindow(UpdateViewModel viewModel)
        {
            InitializeComponent();
            _viewModel = viewModel;
            DataContext = _viewModel;   
        }

        private void Window_Loaded(object sender, RoutedEventArgs e)
        {
            if (_viewModel.AutoStart)
            {
                _viewModel.Message = "Starting...";
                _viewModel.Update.Execute(null);
            }
        }

        private void Window_Closing(object sender, System.ComponentModel.CancelEventArgs e)
        {
            if (_viewModel.IsUpdating)
            {
                var button = MessageBox.Show("Closing will cancel the in-progress update. Close and cancel?", "Confirm close", MessageBoxButton.YesNo);
                if (button == MessageBoxResult.Yes)
                {
                    _viewModel.Cancel();
                }
                else
                {
                    e.Cancel = true;
                }
            }
        }

        private void Window_KeyDown(object sender, KeyEventArgs e)
        {
            switch (e.Key)
            {
                case Key.Escape:
                    Close();
                    return;
            }
        }
    }

    public class UpdateViewModel : ViewModel
    {
        public UpdateViewModel(Version githubVersion, bool autoStart)
        {
            _githubVersion = githubVersion;
            _message = "";
            _autoStart = autoStart;
        }

        private Version _githubVersion;
        public Version Version
        {
            get { return _githubVersion; }
            set
            {
                _githubVersion = value;
                FirePropertyChanged();
                FirePropertyChanged(nameof(VersionString));
            }
        }
        public string VersionString
        {
            get => _githubVersion.ToString();
        }



        private bool _isUpdating;
        public bool IsUpdating
        {
            get { return _isUpdating; }
            set { _isUpdating = value; FirePropertyChanged(); }
        }


        private string _message;
        public string Message
        {
            get { return _message; }
            set { _message = value; FirePropertyChanged(); }
        }


        private bool _autoStart;
        public bool AutoStart
        {
            get { return _autoStart; }
            set
            {
                _autoStart = value;
                FirePropertyChanged();
                FirePropertyChanged(nameof(UpdateButtonVisibility));
            }
        }
        public Visibility UpdateButtonVisibility
        {
            get => _autoStart ? Visibility.Collapsed : Visibility.Visible;
        }

        private CancellationTokenSource? _cts;
        public void Cancel()
        {
            if (_cts != null)
            {
                _cts.Cancel();
                _cts = null;
            }
        }
        public DelegateCommand<object?> Update => new DelegateCommand<object?>(foo =>
        {
            // IIFE style approach to async execution without awaiting
            var _ = ((Func<Task>)(async () =>
            {
                IsUpdating = true;
                try
                {
                    _cts = new CancellationTokenSource();
                    await Updater.UpdateAsync(_cts.Token, status => AppendStatus(status));
                }
                catch (Exception ex)
                {
                    AppendStatus($"Error: {ex.Message}");
                }
                // Not setting IsUpdating to false as needs logic to check if update can be run again
            }))();
        }, _ => !IsUpdating);

        private void AppendStatus(string message)
        {
            Message += "\n" + message;
        }
    }

    public class DesignTimeUpdateViewModel : UpdateViewModel
    {
        public DesignTimeUpdateViewModel()
            : base(new Version("1.2.3"), autoStart: false)
        {
            Message = "Test Message...\nMore here...";
        }
    }
}
