package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	err := validateDomainName("www.flatbo.at", "www")
	assert.NotNil(t, err, "must be fully qualified")

	err = validateDomainName("www.flatbo.at.", "www")
	assert.NotNil(t, err, "www is invalid")

	err = validateDomainName("test.a.b.www.flatbo.at.", "www")
	assert.NotNil(t, err, "www is invalid")

	err = validateDomainName("asdf.flatbo.at.asdf.flatbo.at.", "asdf")
	assert.NotNil(t, err, "messwithdns occurs twice")

	err = validateDomainName("x..flatbo.at.", "asdf")
	assert.NotNil(t, err, "invalid domain name")

	err = validateDomainName("asdf.test.flatbo.at.", "test")
	assert.Nil(t, err)

	err = validateDomainName("a.b.c.d.flatbo.at.", "d")
	assert.Nil(t, err)
}
