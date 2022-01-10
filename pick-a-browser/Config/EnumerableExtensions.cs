using System.Collections.Generic;

namespace pick_a_browser.Config
{
    public static class EnumerableExtensions
    {
        public static IEnumerable<T> NonNulls<T>(this IEnumerable<T?> source)
        {
            foreach(T? item in source)
            {
                if (item != null)
                {
                    yield return item;
                }
            }
        }
    }

}
