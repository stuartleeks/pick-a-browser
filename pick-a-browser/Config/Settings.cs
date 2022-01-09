using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.IO;
using System.Net.NetworkInformation;
using System.Text;
using System.Text.Json;
using System.Text.Json.Nodes;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;

namespace pick_a_browser.Config
{
    public class Settings
    {
        public Settings(Browsers browsers, List<Rule> rules)
        {
            Browsers = browsers;
            Rules = rules;
        }
        public Browsers Browsers { get; }
        public List<Rule> Rules { get; }

        // TODO - add tranformation rules (link shorteners, regex rules)
        // TODO - add browser matches (simple, prefix, regex)

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

            var rules = ParseRules(rootNode);

            return new Settings(browsers, rules);
        }

        private static List<Rule> ParseRules(JsonNode rootNode)
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

        private static Rule ParseRule(JsonNode? ruleNode)
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

        private static Browsers ParseBrowsers(JsonNode? rootNode)
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

        private static Browser ParseBrowser(JsonNode? browserNode)
        {
            if (browserNode == null)
                throw new Exception("browser array item was null");

            var name = browserNode.GetRequiredString("name");
            var exe = browserNode.GetRequiredString("exe");
            var args = browserNode.GetOptionalString("args");
            var iconPath = browserNode.GetOptionalString("iconPath");

            return new Browser(name, exe, args, iconPath);
        }
    }
}
