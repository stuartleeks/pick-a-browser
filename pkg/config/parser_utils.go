package config

import (
	"fmt"
)

func getObjectChildNode(node map[string]interface{}, propertyName string, required bool) (map[string]interface{}, error) {
	if childNode, ok := node[propertyName]; ok {
		switch childNode := childNode.(type) {
		case map[string]interface{}:
			return childNode, nil
		default:
			return map[string]interface{}{}, fmt.Errorf("%q property should be an object", propertyName)
		}
	}
	if required {
		return map[string]interface{}{}, fmt.Errorf("%q property not found", propertyName)
	}
	return map[string]interface{}{}, nil
}

func getArrayChildNode(node map[string]interface{}, propertyName string, required bool) ([]interface{}, error) {
	if childNode, ok := node[propertyName]; ok {
		switch childNode := childNode.(type) {
		case []interface{}:
			return childNode, nil
		default:
			return []interface{}{}, fmt.Errorf("%q property should be an array", propertyName)
		}
	}
	if required {
		return []interface{}{}, fmt.Errorf("%q property not found", propertyName)
	}
	return []interface{}{}, nil
}

func getRequiredString(node map[string]interface{}, propertyName string) (string, error) {
	propertyNode, ok := node[propertyName]
	if !ok {
		return "", fmt.Errorf("required property %q not found", propertyName)
	}
	propertyValue, ok := propertyNode.(string)
	if !ok {
		return "", fmt.Errorf("required property %q expected to be a string", propertyName)
	}
	return propertyValue, nil
}

func getOptionalString(node map[string]interface{}, propertyName string) (*string, error) {
	propertyNode, ok := node[propertyName]
	if !ok {
		return nil, nil
	}
	if propertyNode == nil {
		return nil, nil
	}
	propertyValue, ok := propertyNode.(string)
	if !ok {
		return nil, fmt.Errorf("optional property %q expected to be a string", propertyName)
	}
	return &propertyValue, nil
}

func getOptionalBoolWithDefault(node map[string]interface{}, propertyName string, defaultValue bool) (bool, error) {
	propertyNode, ok := node[propertyName]
	if !ok {
		return defaultValue, nil
	}
	propertyValue, ok := propertyNode.(bool)
	if !ok {
		return false, fmt.Errorf("optional property %q expected to be a bool", propertyName)
	}
	return propertyValue, nil
}
