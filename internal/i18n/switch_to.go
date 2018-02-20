package i18n

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/toby3d/MyPackBot/internal/models"
)

// SwitchTo try load locales by input language codes and return TranslateFunc
func SwitchTo(codes ...string) (T i18n.TranslateFunc, err error) {
	codes = append(codes, models.LanguageFallback)
	T, err = i18n.Tfunc(codes[0], codes[1:]...)
	return
}
