using System;

namespace pick_a_browser.Config
{
    public abstract class Rule
    {
        public const string BrowserPrompt = "_prompt_";

        protected Rule(string browserId)
        {
            BrowserId = browserId;
        }
        public string BrowserId { get; }


        /// <summary>
        /// Test whether the rule matches the specified Uri
        /// </summary>
        /// <param name="uri"></param>
        /// <returns></returns>
        public RuleMatch GetMatch(Uri uri)
        {
            var weight = GetMatchWeight(uri);

            if (weight == 0)
                return RuleMatch.None;

            if (BrowserId == BrowserPrompt)
                weight = int.MaxValue;

            return new RuleMatch(weight, BrowserId);
        }
        protected abstract int GetMatchWeight(Uri uri);
    }

    /// <summary>
    /// Matches a URI using prefix matching
    /// </summary>
    public class PrefixRule : Rule
    {
        public PrefixRule(string prefixMatch, string browserId)
            : base(browserId)
        {
            PrefixMatch = prefixMatch.ToLowerInvariant();
        }

        public string PrefixMatch { get; }

        protected override int GetMatchWeight(Uri uri)
        {
            if (uri.ToString().ToLowerInvariant().StartsWith(PrefixMatch))
                return PrefixMatch.Length; // use prefix length as the weighting
            else
                return 0;
        }
    }

    /// <summary>
    /// Matches a URI based on host suffix
    /// </summary>
    public class HostRule : Rule
    {
        public HostRule(string host, string browserId)
            : base(browserId)
        {
            Host = host.ToLowerInvariant();
        }

        public string Host { get; }

        protected override int GetMatchWeight(Uri uri)
        {
            if (uri.Host.ToLowerInvariant().EndsWith(Host))
                return Host.Length;
            else
                return 0;
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
