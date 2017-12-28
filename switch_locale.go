package main

import (
	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
)

const langFallback = "en"

func switchLocale(langCode string) (T i18n.TranslateFunc, err error) {
	log.Ln("Check", langCode, "localization")
	T, err = i18n.Tfunc(langCode)
	if err != nil {
		log.Ln("Unsupported language, change to ", langFallback, " by default")
		T, err = i18n.Tfunc(langFallback)
	}
	return
}
