package validator

import (
	"encoding/json"
	"strings"
)

// ValidateValue returns true if value is contained in allowList
func ValidateValue(value string, allowList string, separator string) bool {
	parts, err := parseValues(allowList, separator)
	if err != nil {
		return false
	}
	for _, part := range parts {
		if part == value {
			return true
		}
	}
	return false
}

func parseValues(allowList, separator string) ([]string, error) {
	if separator == "json" {
		var parts []string
		err := json.Unmarshal([]byte(allowList), &parts)
		return parts, err
	}
	return strings.Split(allowList, separator), nil
}

// ValidateCommonName returns true if value is contained in allowList
func ValidateCommonName(value string, allowList string, separator string) bool {
	return ValidateValue(value, allowList, separator)
}
