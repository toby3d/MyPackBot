package i18n

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"gitlab.com/toby3d/mypackbot/internal/models"
)

// SwitchTo try load locales by input language codes and return TranslateFunc
func SwitchTo(codes ...string) (t i18n.TranslateFunc, err error) {
	codes = append(codes, models.LanguageFallback)
	t, err = i18n.Tfunc(codes[0], codes[1:]...)
	return
}
