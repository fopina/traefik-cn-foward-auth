package validator

import (
	"encoding/json"
	"net/url"
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
	parts, err := parseValues(allowList, separator)
	if err != nil {
		return false
	}
	decodedStr, err := url.QueryUnescape(value)
	if err != nil {
		return false
	}
	// This is the hard part: traefik does not properly serialize these fields: https://github.com/traefik/traefik/issues/11309
	// best effort and fallback to invalid
	cnParts := strings.Split(decodedStr, ",")
	// as modest attempt to assert first part is full, check 2nd part expected prefix
	if (len(cnParts) > 1) && (!strings.HasPrefix(cnParts[1], "Subject=\"CN=")) {
		return false
	}
	if !strings.HasPrefix(cnParts[0], "Subject=\"CN=") {
		return false
	}
	// traefik sends full chain of certificates sent by the client - we only care about the leaf one (first) - rest of chain is validated by the mTLS setup
	finalPart := cnParts[0][12 : len(cnParts[0])-1]
	for _, part := range parts {
		if part == finalPart {
			return true
		}
	}
	return false
}
