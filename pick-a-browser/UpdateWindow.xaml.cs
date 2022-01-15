using pick_a_browser.Config;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;
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
        public UpdateWindow(UpdateViewModel viewModel)
        {
            InitializeComponent();
            DataContext = viewModel;
        }
    }

    public class UpdateViewModel : ViewModel
    {
        public UpdateViewModel(Version githubVersion)
        {
            _githubVersion = githubVersion;
            _message = "";
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


        public DelegateCommand<object?> Update => new DelegateCommand<object?>(foo =>
        {
            // IIFE style approach to async execution without awaiting
            var _ = ((Func<Task>)(async () =>
            {
                IsUpdating = true;
                try
                {
                    await Updater.UpdateAsync(status => AppendStatus(status));
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
            : base(new Version("1.2.3"))
        {
            Message = "Test Message...\nMore here...";
        }
    }
}
