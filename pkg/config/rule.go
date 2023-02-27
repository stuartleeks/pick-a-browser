package config

import (
	"fmt"
	"net/url"
	"strings"
)

type Rule interface {
	Type() string
	BrowserId() string
	Match(url *url.URL) int
}

type RuleBase struct {
	browserId string
}

func (r *RuleBase) BrowserId() string {
	return r.browserId
}

type PrefixRule struct {
	RuleBase
	prefixMatch string
}

func (r *PrefixRule) Type() string {
	return "prefix"
}

func (r *PrefixRule) Match(url *url.URL) int {
	if strings.HasPrefix(url.String(), r.prefixMatch) {
		return len(r.prefixMatch)
	}
	return -1
}

var _ Rule = &PrefixRule{}

type HostRule struct {
	RuleBase
	host string
}

func (r *HostRule) Type() string {
	return "host"
}

func (r *HostRule) Match(url *url.URL) int {
	if strings.HasSuffix(url.Host, r.host) {
		return len(r.host)
	}
	return -1
}

var _ Rule = &HostRule{}

func parseRules(rootNode map[string]interface{}) ([]Rule, error) {
	rulesNode, err := getArrayChildNode(rootNode, "rules", false)
	if err != nil {
		return []Rule{}, err
	}

	rules := []Rule{}

	for _, v := range rulesNode {
		ruleNode, ok := v.(map[string]interface{})
		if !ok {
			return []Rule{}, fmt.Errorf("expected objects in rules array")
		}
		rule, err := parseRule(ruleNode)
		if err != nil {
			return []Rule{}, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func parseRule(ruleNode map[string]interface{}) (Rule, error) {

	ruleType, err := getRequiredString(ruleNode, "type")
	if err != nil {
		return nil, err
	}
	browser, err := getRequiredString(ruleNode, "browser")
	if err != nil {
		return nil, err
	}

	switch ruleType {
	case "prefix":
		prefix, err := getRequiredString(ruleNode, "prefix")
		if err != nil {
			return nil, err
		}
		return &PrefixRule{
			prefixMatch: prefix,
			RuleBase:    RuleBase{browserId: browser},
		}, nil
	case "host":
		host, err := getRequiredString(ruleNode, "host")
		if err != nil {
			return nil, err
		}
		return &HostRule{
			host:     host,
			RuleBase: RuleBase{browserId: browser},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported rule type %q", ruleType)
	}
}

func MatchRules(rules []Rule, url *url.URL) string {
	matchWeight := -1
	browserId := ""
	for _, rule := range rules {
		tmpWeight := rule.Match(url)
		if tmpWeight > matchWeight {
			matchWeight = tmpWeight
			browserId = rule.BrowserId()
		}
	}
	return browserId
}
