package middleware

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetPrinter(code string) *message.Printer {
	return message.NewPrinter(language.Make(code))
}
