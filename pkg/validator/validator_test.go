package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValue(t *testing.T) {
	assert.True(t, ValidateValue("ehlo", "helo,ehlo", ","))
	assert.False(t, ValidateValue("ehlo", "helo", ","))
	assert.True(t, ValidateValue("ehlo", "helo|ehlo", "|"))
}

func TestValidateValueJSON(t *testing.T) {
	assert.True(t, ValidateValue("ehlo", `["helo","ehlo"]`, "json"))
	assert.False(t, ValidateValue("ehlo", `["helo"]`, "json"))
	// invalid json
	assert.False(t, ValidateValue("ehlo", `["helo","ehlo"`, "json"))
}
