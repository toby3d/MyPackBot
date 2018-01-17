package main

import (
	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
)

const langFallback = "en"

func switchLocale(langCode string) (T i18n.TranslateFunc, err error) {
	log.Ln("Check", langCode, "localization")
	T, err = i18n.Tfunc(langCode, langFallback)
	errCheck(err)
	return
}
