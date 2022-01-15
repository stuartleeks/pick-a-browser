using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Nodes;
using System.Threading.Tasks;

namespace pick_a_browser.Config
{
    public class SettingsSerialization
    {
        public static async Task<Settings> LoadFromFileAsync(string filename)
        {
            var settingsContent = await File.ReadAllTextAsync(filename);
            return ParseSettings(settingsContent);
        }

        public static Settings ParseSettings(string settingsContent)
        {
            var rootNode = JsonNode.Parse(settingsContent, null, new JsonDocumentOptions { AllowTrailingCommas = true, CommentHandling = JsonCommentHandling.Skip });

            if (rootNode == null)
                throw new Exception("Failed to parse settings");

            var browsers = ParseBrowsers(rootNode);
            var transformations = ParseTransformations(rootNode);
            var rules = ParseRules(rootNode);

            var updateCheck = ParseUpdateCheck(rootNode);

            return new Settings(browsers, transformations, rules, updateCheck);
        }

        public static Browsers ParseBrowsers(JsonNode rootNode)
        {
            var browsersNode = rootNode["browsers"];
            if (browsersNode == null)
                throw new Exception("browsers not found in settings");

            var browsers = new List<Browser>();
            foreach (var browserNode in browsersNode.AsArray())
            {
                var browser = ParseBrowser(browserNode);
                browsers.Add(browser);
            }

            return new Browsers(browsers);
        }
        public static JsonNode ToJsonNode(Browsers browsers)
        {
            return new JsonArray(
                browsers.Select(ToJsonNode).ToArray()
                );
        }

        public static Browser ParseBrowser(JsonNode? browserNode)
        {
            if (browserNode == null)
                throw new Exception("browser array item was null");

            var id= browserNode.GetRequiredString("id");
            var name = browserNode.GetRequiredString("name");
            var exe = browserNode.GetRequiredString("exe");
            var args = browserNode.GetOptionalString("args");
            var iconPath = browserNode.GetOptionalString("iconPath");
            var hidden = (bool?)browserNode["hidden"] ?? false;

            return new Browser(id, name, exe, args, iconPath, hidden);
        }
        public static JsonNode ToJsonNode(Browser browser)
        {
            var node = new JsonObject();
            node.Add("id", browser.Id);
            node.Add("name", browser.Name);
            node.Add("exe", browser.Exe);
            if (browser.Args != null)
                node.Add("args", browser.Args);
            if (browser.IconPath != null)
                node.Add("iconPath", browser.IconPath);
            node.Add("hidden", browser.Hidden);
            return node;
        }
        public static List<Rule> ParseRules(JsonNode rootNode)
        {
            var rules = new List<Rule>();
            var rulesNode = rootNode["rules"];
            if (rulesNode == null)
                return rules;

            foreach (var ruleNode in rulesNode.AsArray())
            {
                var rule = ParseRule(ruleNode);
                rules.Add(rule);
            }

            return rules;
        }

        public static Rule ParseRule(JsonNode? ruleNode)
        {
            if (ruleNode == null)
                throw new Exception("rule array item was null");

            var type = ruleNode.GetRequiredString("type");

            var browser = ruleNode.GetRequiredString("browser");
            switch (type.ToLowerInvariant())
            {
                case "prefix":
                    var prefixMatch = ruleNode.GetRequiredString("prefix");
                    return new PrefixRule(prefixMatch, browser);
                case "host":
                    var host = ruleNode.GetRequiredString("host");
                    return new HostRule(host, browser);
                default:
                    throw new Exception($"Unsupported rule type: '{type}'");
            }
        }

        public static Transformations ParseTransformations(JsonNode rootNode)
        {
            var transformationsNode = rootNode["transformations"];
            if (transformationsNode == null)
                return new Transformations(new List<string>(), new List<LinkWrapper>());

            var linkShorteners = ParseLinkShorteners(transformationsNode);
            var linkWrappers= ParseLinkWrappers(transformationsNode);

            return new Transformations(linkShorteners, linkWrappers);
        }
        public static List<string> ParseLinkShorteners(JsonNode transformationsNode)
        {
            var shorteners = new List<string>();
            var shortenersNode = transformationsNode["linkShorteners"];
            if (shortenersNode == null)
                return shorteners;

            foreach (var shortenerNode in shortenersNode.AsArray())
            {
                var shortener = (string?)shortenerNode;
                if (string.IsNullOrEmpty(shortener))
                {
                    throw new Exception("linkShorteners must not be null or empty");
                }
                shorteners.Add(shortener);
            }

            return shorteners;
        }
        public static List<LinkWrapper> ParseLinkWrappers(JsonNode transformationsNode)
        {
            var linkWrappersNode = transformationsNode["linkWrappers"];
            if (linkWrappersNode == null)
                return new List<LinkWrapper>();

            var linkWrappers = new List<LinkWrapper>();
            foreach (var linkWrapperNode in linkWrappersNode.AsArray())
            {
                var linkWrapper = ParseLinkWrapper(linkWrapperNode);
                linkWrappers.Add(linkWrapper);
            }

            return linkWrappers;
        }
        public static LinkWrapper ParseLinkWrapper(JsonNode? linkWrapperNode)
        {
            if (linkWrapperNode == null)
                throw new Exception("linkWrapper array item was null");

            var urlPrefix = linkWrapperNode.GetRequiredString("prefix");
            var queryStringKey = linkWrapperNode.GetRequiredString("queryString");
            
            return new LinkWrapper(urlPrefix, queryStringKey);
        }
        public static UpdateCheck ParseUpdateCheck(JsonNode rootNode)
        {
            var updatesNode = rootNode["updates"];
            if (updatesNode == null)
                return UpdateCheck.Prompt;

            var updatesString = (string?)updatesNode;
            if (updatesString == null)
                throw new Exception($"updates property must be set to a string value: {string.Join(", ", Enum.GetNames(typeof(UpdateCheck)).Select(s => s.ToLower()))}");

            object? updateCheckResult;
            if (!Enum.TryParse(typeof(UpdateCheck), updatesString, true, out updateCheckResult))
                throw new Exception($"updates property must one of the following values: {string.Join(", ", Enum.GetNames(typeof(UpdateCheck)).Select(s => s.ToLower()))}");

            return (UpdateCheck) updateCheckResult!;
        }

    }
}
