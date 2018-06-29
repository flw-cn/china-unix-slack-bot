package util

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedRegexpParse(t *testing.T) {
	text := "Hello, alice!"
	re := regexp.MustCompile(`^Hello,\s+(?P<name>\w+)`)
	ok, dict := NamedRegexpParse(text, re)
	assert.True(t, ok, "must matched")
	assert.NotNil(t, dict, "must captured")
	assert.NotNil(t, dict["name"], "must captured name")
	assert.Equal(t, "alice", dict["name"], "name not correct")
}
