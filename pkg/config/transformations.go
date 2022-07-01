package config

import (
	"fmt"
)

type Transformations struct {
	LinkShorteners []string
	LinkWrappers   []LinkWrapper
}

type LinkWrapper struct {
	UrlPrefix      string
	QueryStringKey string
}

func GetDefaultLinkShorteners() []string {
	return []string{
		"aka.ms",
		"t.co",
		"go.microsoft.com",
	}
}
func GetDefaultLinkWrappers() []LinkWrapper {
	return []LinkWrapper{
		{
			UrlPrefix:      "https://statics.teams.cdn.office.net/evergreen-assets/safelinks/",
			QueryStringKey: "url",
		},
		{
			UrlPrefix:      "https://staticsint.teams.cdn.office.net/evergreen-assets/safelinks/",
			QueryStringKey: "url",
		},
		{
			UrlPrefix:      "https://nam06.safelinks.protection.outlook.com/",
			QueryStringKey: "url",
		},
	}
}

func parseTransformations(rootNode map[string]interface{}) (Transformations, error) {
	transformationsNode, err := getObjectChildNode(rootNode, "transformations", false)
	if err != nil {
		return Transformations{}, err
	}

	linkShorteners, err := parseLinkShorteners(transformationsNode)
	if err != nil {
		return Transformations{}, err
	}
	linkWrappers, err := parseLinkWrappers(transformationsNode)
	if err != nil {
		return Transformations{}, err
	}

	return Transformations{
		LinkShorteners: linkShorteners,
		LinkWrappers:   linkWrappers,
	}, nil
}
func parseLinkShorteners(parentNode map[string]interface{}) ([]string, error) {
	linkShortenersNode, err := getArrayChildNode(parentNode, "linkShorteners", false)
	if err != nil {
		return []string{}, err
	}

	linkShorteners := []string{}

	for _, v := range linkShortenersNode {
		linkShortener, ok := v.(string)
		if !ok {
			return []string{}, fmt.Errorf("expected string entry in linkShorteners array")
		}
		linkShorteners = append(linkShorteners, linkShortener)
	}

	return linkShorteners, nil
}

func parseLinkWrappers(parentNode map[string]interface{}) ([]LinkWrapper, error) {
	linkWrappersNode, err := getArrayChildNode(parentNode, "linkWrappers", false)
	if err != nil {
		return []LinkWrapper{}, err
	}

	linkWrappers := []LinkWrapper{}

	for _, v := range linkWrappersNode {
		linkWrapperNode, ok := v.(map[string]interface{})
		if !ok {
			return []LinkWrapper{}, fmt.Errorf("expected objects in linkWrappers array")
		}
		linkWrapper, err := parseLinkWrapper(linkWrapperNode)
		if err != nil {
			return []LinkWrapper{}, err
		}
		linkWrappers = append(linkWrappers, linkWrapper)
	}

	return linkWrappers, nil
}
func parseLinkWrapper(linkWrapperNode map[string]interface{}) (LinkWrapper, error) {

	prefix, err := getRequiredString(linkWrapperNode, "prefix")
	if err != nil {
		return LinkWrapper{}, err
	}
	queryStringKey, err := getRequiredString(linkWrapperNode, "queryString")
	if err != nil {
		return LinkWrapper{}, err
	}

	return LinkWrapper{
		UrlPrefix:      prefix,
		QueryStringKey: queryStringKey,
	}, nil
}
