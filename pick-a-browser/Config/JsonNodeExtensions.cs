using System;
using System.Text.Json.Nodes;

namespace pick_a_browser.Config
{
    public static class JsonNodeExtensions
    {
        public static string GetRequiredString(this JsonNode node, string name)
        {
            var value = node[name];
            if (value == null)
                throw new Exception($"Property '{name}' not found");

            var stringValue = (string?)value;
            if (stringValue == null)
                throw new Exception($"Property '{name}' is null");

            return stringValue;
        }
        public static string? GetOptionalString(this JsonNode node, string name)
        {
            var value = node[name];
            if (value == null)
                return null;
            else
                return (string?)value;

        }

    }
}
