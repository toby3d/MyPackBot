package utils

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

//nolint:gochecknoglobals
var (
	// Skin colors for remove
	bannedSkins = []rune{'🏻', '🏼', '🏽', '🏾', '🏿'}

	// Transformer for remove skin colors
	skinRemover = runes.Remove(runes.Predicate(func(r rune) bool {
		for _, skin := range bannedSkins {
			if r != skin {
				continue
			}
			return true
		}
		return false
	}))
)

// FixEmojiTone remove any skin tone from input emoji.
func FixEmojiTone(raw string) (string, error) {
	result, _, err := transform.String(skinRemover, raw)
	return result, err
}
