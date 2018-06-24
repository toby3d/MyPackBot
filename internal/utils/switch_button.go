package utils

import (
	"github.com/nicksnyder/go-i18n/i18n"
	tg "gitlab.com/toby3d/telegram"
)

// SwitchButton helper just generate ReplyMarkup with SelfSwitch button
func SwitchButton(T i18n.TranslateFunc) *tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitchSelf(T("button_inline_select"), " "),
		),
	)
}
