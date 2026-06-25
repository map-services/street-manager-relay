package cmd

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestShouldExclude(t *testing.T) {
	excludeTypes := []string{"SkipS", "hoarding"}

	excluded := make(map[string]bool)
	for _, t := range excludeTypes {
		excluded[strings.ToLower(t)] = true
	}

	assert.True(t, excluded[strings.ToLower("Skips")])
	assert.True(t, excluded[strings.ToLower("Hoarding")])
	assert.False(t, excluded[strings.ToLower("Utility repair")])
}
