package utils

import (
	"github.com/nicksnyder/go-i18n/i18n"
	tg "gitlab.com/toby3d/telegram"
)

// MenuKeyboard helper just generate ReplyMarkup with menu buttons
func MenuKeyboard(t i18n.TranslateFunc) (rkm *tg.ReplyKeyboardMarkup) {
	rkm = tg.NewReplyKeyboardMarkup(
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(t("button_add_sticker")),
			tg.NewReplyKeyboardButton(t("button_add_pack")),
		),
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(t("button_del_sticker")),
			tg.NewReplyKeyboardButton(t("button_del_pack")),
		),
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(t("button_reset")),
		),
	)

	rkm.ResizeKeyboard = true
	return
}
