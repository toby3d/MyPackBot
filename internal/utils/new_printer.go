package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NewPrinter(lang string) *message.Printer {
	return message.NewPrinter(message.MatchLanguage(
		lang,
		language.English.String(),
	))
}
