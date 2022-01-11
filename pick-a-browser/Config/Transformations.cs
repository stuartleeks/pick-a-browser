using System.Collections.Generic;

namespace pick_a_browser.Config
{
    public class Transformations
    {
        public Transformations(List<string> linkShorteners, List<LinkWrapper> linkWrappers)
        {
            LinkShorteners = linkShorteners;
            LinkWrappers = linkWrappers;
        }

        public List<string> LinkShorteners { get; }
        public List<LinkWrapper> LinkWrappers { get; }
    }

    public class LinkWrapper
    {
        public LinkWrapper(string urlPrefix, string queryStringKey)
        {
            UrlPrefix = urlPrefix;
            QueryStringKey = queryStringKey;
        }

        public string UrlPrefix { get; }
        public string QueryStringKey { get; }
    }
}
