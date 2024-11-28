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

func TestValidateCommonName(t *testing.T) {
	assert.True(t, ValidateCommonName("Subject%3D%22CN%3Dauth-client%22", `auth-client`, ","))
	assert.True(t, ValidateCommonName("Subject%3D%22CN%3Dauth-client%22%2CSubject%3D%22CN%3Dgood-one%22", `auth-client`, ","))
	assert.False(t, ValidateCommonName("Subject%3D%22CN%3Dauth-client%22%2CSubject%3D%22CN%3Dgood-one%22", `not-auth-client`, ","))
}

func TestValidateCommonNameJSON(t *testing.T) {
	assert.True(t, ValidateCommonName("Subject%3D%22CN%3Dauth-client%22", `["auth-client"]`, "json"))
}
