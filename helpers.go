package main

import (
	"github.com/nicksnyder/go-i18n/i18n"
	tg "github.com/toby3d/telegram"
)

func getMenuKeyboard(T i18n.TranslateFunc) *tg.ReplyKeyboardMarkup {
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

func getCancelButton(T i18n.TranslateFunc) *tg.ReplyKeyboardMarkup {
	return tg.NewReplyKeyboardMarkup(
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(T("button_cancel")),
		),
	)
}
