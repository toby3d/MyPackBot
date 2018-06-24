package utils

import (
	"github.com/nicksnyder/go-i18n/i18n"
	tg "gitlab.com/toby3d/telegram"
)

// MenuKeyboard helper just generate ReplyMarkup with menu buttons
func MenuKeyboard(T i18n.TranslateFunc) *tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboardMarkup(
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(T("button_add_sticker")),
			tg.NewReplyKeyboardButton(T("button_add_pack")),
		),
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(T("button_del_sticker")),
			tg.NewReplyKeyboardButton(T("button_del_pack")),
		),
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(T("button_reset")),
		),
	)
}
