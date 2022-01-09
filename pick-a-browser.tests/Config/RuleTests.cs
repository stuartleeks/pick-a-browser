using pick_a_browser.Config;
using System;
using Xunit;

namespace pick_a_browser.tests.Config
{
    public class RuleTests
    {
        [Fact]
        public void PrefixRule_MatchOnSharedPrefix()
        {
            var rule = new PrefixRule("https://example.com", "browser");
            var match = rule.GetMatch(new Uri("https://example.com/foo"));
            Assert.NotNull(match);
            Assert.True(match.Weight > 0);
        }
        [Fact]
        public void PrefixRule_NoMatchOnDifferentPrefix()
        {
            var rule = new PrefixRule("https://example.com", "browser");
            var match = rule.GetMatch(new Uri("https://www.example.com/foo"));
            Assert.NotNull(match);
            Assert.True(match.Weight == 0);
        }


        [Fact]
        public void HostRule_MatchOnExactHostMatch()
        {
            var rule = new HostRule("example.com", "browser");
            var match = rule.GetMatch(new Uri("https://example.com/foo"));
            Assert.NotNull(match);
            Assert.True(match.Weight > 0);
        }
        [Fact]
        public void HostRule_MatchOnSharedSuffix()
        {
            var rule = new HostRule("example.com", "browser");
            var match = rule.GetMatch(new Uri("https://foo.example.com/foo"));
            Assert.NotNull(match);
            Assert.True(match.Weight > 0);
        }

        [Fact]
        public void HostRule_NoMatchOnDifferentHost()
        {
            var rule = new PrefixRule("https://example.com", "browser");
            var match = rule.GetMatch(new Uri("https://www.example.net/foo"));
            Assert.NotNull(match);
            Assert.True(match.Weight == 0);
        }


    }
}