package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixEmojiTone(t *testing.T) {
	for in, out := range map[string]string{
		"👍":     "👍",
		"✌🏻":    "✌",
		"👋🏽👌🤙🏿": "👋👌🤙",
		"":      "",
	} {
		in, out := in, out
		t.Run(in, func(t *testing.T) {
			result, err := FixEmojiTone(in)
			assert.NoError(t, err)
			assert.Equal(t, out, result)
		})
	}
}
