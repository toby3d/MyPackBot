package utils

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// Skin colors for remove
var bannedSkins = []rune{127995, 127996, 127997, 127998, 127999}

// Transformer for remove skin colors
var skinRemover = runes.Remove(runes.Predicate(
	func(r rune) bool {
		for _, skin := range bannedSkins {
			if r == skin {
				return true
			}
		}
		return false
	},
))

// FixEmoji fixes user input by remove all potential skin colors
func FixEmoji(raw string) (string, error) {
	result, _, err := transform.String(skinRemover, raw)
	return result, err
}
