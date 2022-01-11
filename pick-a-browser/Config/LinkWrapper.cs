namespace pick_a_browser.Config
{
    public class LinkWrapper
    {
        public LinkWrapper(string urlPrefix, string queryStringKey)
        {
            UrlPrefix = urlPrefix;
            QueryStringKey = queryStringKey;
        }

        public string UrlPrefix { get; }
        public string QueryStringKey { get;  }
    }
}
