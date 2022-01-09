using System.Collections;
using System.Collections.Generic;
using System.Linq;

namespace pick_a_browser.Config
{
    /// <summary>
    /// Represents a readonly list of Browser instances allowing for indexed access on the Name property
    /// </summary>
    public class Browsers : IEnumerable<Browser>, IReadOnlyList<Browser>
    {
        private readonly List<Browser> _browsers;

        public Browsers(List<Browser> browsers)
        {
            _browsers = browsers;
            // TODO - ensure that names are unique?
        }

        public Browser this[int index] => _browsers[index];
        public Browser this[string name] => _browsers.First(b=>b.Name == name);

        public int Count => _browsers.Count;

        public IEnumerator<Browser> GetEnumerator()
        {
            return _browsers.GetEnumerator();
        }

        IEnumerator IEnumerable.GetEnumerator()
        {
            return GetEnumerator();
        }
    }

    public class Browser
    {
        public Browser(string name, string exe, string? args, string? iconPath)
        {
            Name = name;
            Exe = exe;
            Args = args;
            IconPath = iconPath;
        }
        public string Name { get; }
        public string Exe { get; }
        public string? Args { get; }
        public string? IconPath { get; }
    }

}
