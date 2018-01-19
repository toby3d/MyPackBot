package main

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

const langFallback = "en"

func switchLocale(langCode string) (i18n.TranslateFunc, error) {
	log.Ln("Check", langCode, "localization")
	return i18n.Tfunc(langCode, langFallback)
}
