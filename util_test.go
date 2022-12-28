package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	assert.Equal(t, ExtractSubdomain("a.b.flatbo.at."), "b")
	assert.Equal(t, ExtractSubdomain("test.com"), "")
}
