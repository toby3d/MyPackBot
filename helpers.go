package main

import (
	"github.com/nicksnyder/go-i18n/i18n"
	tg "github.com/toby3d/telegram"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var bannedSkins = []rune{127995, 127996, 127997, 127998, 127999}

var skinRemover = runes.Remove(runes.Predicate(
	func(r rune) bool {
		for _, skin := range bannedSkins {
			if r == skin {
				return true
			}
		}
		return false
	},
))

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

func getSwitchButton(T i18n.TranslateFunc) *tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitchSelf(T("button_inline_select"), " "),
		),
	)
}

func fixEmoji(raw string) (string, error) {
	result, _, err := transform.String(skinRemover, raw)
	return result, err
}
