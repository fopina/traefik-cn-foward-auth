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

func TestValidateValueAllowedJSON(t *testing.T) {
	assert.True(t, validateValueAllowedJSON("ehlo", `["helo","ehlo"]`))
	assert.True(t, ValidateValue("ehlo", `["helo","ehlo"]`, "json"))
	assert.False(t, validateValueAllowedJSON("ehlo", `["helo"]`))
	// invalid json
	assert.False(t, ValidateValue("ehlo", `["helo","ehlo"`, "json"))
}
