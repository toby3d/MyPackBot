package utils

import (
	"github.com/nicksnyder/go-i18n/i18n"
	tg "gitlab.com/toby3d/telegram"
)

// CancelButton helper just generate ReplyMarkup with cancel button
func CancelButton(T i18n.TranslateFunc) *tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboardMarkup(
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(T("button_cancel")),
		),
	)
}
