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
            ""id"" : ""test1"",
            ""name"" : ""test1"",
            ""exe"" : ""test1_exe"",
        },
        {
            ""id"" : ""test2"",
            ""name"" : ""test2"",
            ""exe"" : ""test2_exe"",
            ""args"" : ""arg1 arg2"",
            ""iconPath"" : ""c:\\some\\path"",
            ""hidden"" : true,
        }
    ]
}";

            var settings = SettingsSerialization.ParseSettings(json);

            Assert.NotNull(settings.Browsers);
            Assert.Equal(2, settings.Browsers.Count);

            var browser = settings.Browsers[0];
            Assert.NotNull(browser);
            Assert.Equal("test1", browser.Id);
            Assert.Equal("test1", browser.Name);
            Assert.Equal("test1_exe", browser.Exe);
            Assert.Null(browser.Args);
            Assert.Null(browser.IconPath);
            Assert.False(browser.Hidden);

            browser = settings.Browsers[1];
            Assert.NotNull(browser);
            Assert.Equal("test2", browser.Id);
            Assert.Equal("test2", browser.Name);
            Assert.Equal("test2_exe", browser.Exe);
            Assert.Equal("arg1 arg2", browser.Args);
            Assert.Equal("c:\\some\\path", browser.IconPath);
            Assert.True(browser.Hidden);
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
            Assert.Equal("https://example.com", prefixRule!.PrefixMatch);
            Assert.Equal("browser1", prefixRule.BrowserId);

            var hostRule = settings.Rules[1] as HostRule;
            Assert.NotNull(hostRule);
            Assert.Equal("example.com", hostRule!.Host);
            Assert.Equal("browser2", hostRule.BrowserId);
        }


        [Fact]
        public void ParseLinkShorteners()
        {
            var json = @"{
    ""browsers"": [],
    ""transformations"" : {
        ""linkShorteners"": [
            ""shortener1"",
            ""shortener2"",
        ]
    },
    ""rules"": []
}";

            var settings = SettingsSerialization.ParseSettings(json);

            Assert.NotNull(settings.Transformations);
            Assert.NotNull(settings.Transformations.LinkShorteners);
            Assert.Equal(2, settings.Transformations.LinkShorteners.Count);

            Assert.Equal("shortener1", settings.Transformations.LinkShorteners[0]);
            Assert.Equal("shortener2", settings.Transformations.LinkShorteners[1]);
        }

        [Fact]
        public void ParseLinkWrappers()
        {
            var json = @"{
    ""browsers"": [],
    ""transformations"" : {
        ""linkWrappers"": [
            { ""prefix"": ""https://example.com"", ""queryString"" : ""url"" },
            { ""prefix"": ""https://example.net"", ""queryString"" : ""u"" },
        ]
    },
    ""rules"": []
}";

            var settings = SettingsSerialization.ParseSettings(json);

            Assert.NotNull(settings.Transformations);
            Assert.NotNull(settings.Transformations.LinkWrappers);
            Assert.Equal(2, settings.Transformations.LinkWrappers.Count);

            var wrapper = settings.Transformations.LinkWrappers[0];
            Assert.Equal("https://example.com", wrapper.UrlPrefix);
            Assert.Equal("url", wrapper.QueryStringKey);

            wrapper = settings.Transformations.LinkWrappers[1];
            Assert.Equal("https://example.net", wrapper.UrlPrefix);
            Assert.Equal("u", wrapper.QueryStringKey);
        }
    }
}