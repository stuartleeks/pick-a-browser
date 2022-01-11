using System;
using System.Windows.Input;

namespace pick_a_browser
{
    public class DelegateCommand<T> : ICommand
    {
        public event EventHandler? CanExecuteChanged;

        private readonly Action<T?> _execute;
        private readonly Predicate<T?>? _canExecute;

        public DelegateCommand(Action<T?> execute)
                       : this(execute, null)
        {
        }
        public DelegateCommand(Action<T?> execute,
                       Predicate<T?>? canExecute)
        {
            _execute = execute;
            _canExecute = canExecute;
        }

        public bool CanExecute(object? parameter) => _canExecute?.Invoke((T?)parameter) ?? true;
        public void Execute(object? parameter) => _execute((T?)parameter);
        public void RaiseCanExecuteChanged() => CanExecuteChanged?.Invoke(this, EventArgs.Empty);
    }

}
