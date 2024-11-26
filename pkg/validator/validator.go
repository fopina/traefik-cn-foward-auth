package validator

import (
	"encoding/json"
	"strings"
)

// ValidateValue returns true if value is contained in allowList
func ValidateValue(value string, allowList string, separator string) bool {
	if separator == "json" {
		return validateValueAllowedJSON(value, allowList)
	}
	parts := strings.Split(allowList, separator)
	for _, part := range parts {
		if part == value {
			return true
		}
	}
	return false
}

func validateValueAllowedJSON(value string, allowList string) bool {
	var parts []string
	err := json.Unmarshal([]byte(allowList), &parts)
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
