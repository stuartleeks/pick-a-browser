using System;

namespace pick_a_browser.Config
{
    public abstract class Rule
    {

        /// <summary>
        /// Test whether the rule matches the specified Uri
        /// </summary>
        /// <param name="uri"></param>
        /// <returns></returns>
        public abstract RuleMatch GetMatch(Uri uri);
    }

    /// <summary>
    /// Matches a URI using prefix matching
    /// </summary>
    public class PrefixRule : Rule
    {
        public PrefixRule(string prefixMatch, string browser)
        {
            PrefixMatch = prefixMatch.ToLowerInvariant();
            Browser = browser;
        }

        public string PrefixMatch { get; }
        public string Browser { get; }

        public override RuleMatch GetMatch(Uri uri)
        {
            if (uri.ToString().ToLowerInvariant().StartsWith(PrefixMatch))
                return new RuleMatch(PrefixMatch.Length, Browser); // use prefix length as the weighting
            else
                return RuleMatch.None;
        }
    }

    /// <summary>
    /// Matches a URI based on host suffix
    /// </summary>
    public class HostRule : Rule
    {
        public HostRule(string host, string browser)
        {
            Host = host.ToLowerInvariant();
            Browser = browser;
        }

        public string Host { get; }
        public string Browser { get; }

        public override RuleMatch GetMatch(Uri uri)
        {
            if (uri.Host.ToLowerInvariant().EndsWith(Host))
                return new RuleMatch(Host.Length, Browser);
            else
                return RuleMatch.None;
        }
    }

    public class RuleMatch
    {
        public static readonly RuleMatch None = new RuleMatch();

        private RuleMatch()
        {
            Weight = 0;
            BrowserId = "";
        }
        public RuleMatch(int weight, string browserId)
        {
            if (weight <= 0)
                throw new ArgumentOutOfRangeException("weight must be greater than zero (or use RuleMatch.None)");
            if (string.IsNullOrEmpty(browserId))
                throw new ArgumentNullException("browserId must be set");

            Weight = weight;
            BrowserId = browserId;
        }

        /// <summary>
        /// The match weight. Zero indicates no match. Positive value indicates a match, with higher number results taking precedence when there are multiple matches
        /// </summary>
        public int Weight { get; }
        public string BrowserId { get; }
    }
}
