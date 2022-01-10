using pick_a_browser.Config;
using Xunit;

namespace pick_a_browser.tests.Config
{
    public class SettingsSerializationTests
    {
        [Fact]
        public void ParseBrowsers()
        {
            var json = @"{
    ""browsers"": [
        {
            ""name"" : ""test1"",
            ""exe"" : ""test1_exe"",
        },
        {
            ""name"" : ""test2"",
            ""exe"" : ""test2_exe"",
            ""args"" : ""arg1 arg2"",
            ""iconPath"" : ""c:\\some\\path"",
        }
    ]
}";

            var settings = SettingsSerialization.ParseSettings(json);

            Assert.NotNull(settings.Browsers);
            Assert.Equal(2, settings.Browsers.Count);

            var browser = settings.Browsers[0];
            Assert.NotNull(browser);
            Assert.Equal("test1", browser.Name);
            Assert.Equal("test1_exe", browser.Exe);
            Assert.Null(browser.Args);
            Assert.Null(browser.IconPath);

            browser = settings.Browsers[1];
            Assert.NotNull(browser);
            Assert.Equal("test2", browser.Name);
            Assert.Equal("test2_exe", browser.Exe);
            Assert.Equal("arg1 arg2", browser.Args);
            Assert.Equal("c:\\some\\path", browser.IconPath);
        }

        [Fact]
        public void ParseRules()
        {
            var json = @"{
    ""browsers"": [],
    ""rules"": [
        {
            ""type"" : ""prefix"",
            ""prefix"" : ""https://example.com"",
            ""browser"" : ""browser1"",
        },
        {
            ""type"" : ""host"",
            ""host"" : ""example.com"",
            ""browser"" : ""browser2"",
        },
    ]
}";

            var settings = SettingsSerialization.ParseSettings(json);

            Assert.NotNull(settings.Rules);
            Assert.Equal(2, settings.Rules.Count);

            var prefixRule = settings.Rules[0] as PrefixRule;
            Assert.NotNull(prefixRule);
            Assert.Equal("https://example.com", prefixRule.PrefixMatch);
            Assert.Equal("browser1", prefixRule.Browser);

            var hostRule = settings.Rules[1] as HostRule;
            Assert.NotNull(hostRule);
            Assert.Equal("example.com", hostRule.Host);
            Assert.Equal("browser2", hostRule.Browser);
        }
    }
}