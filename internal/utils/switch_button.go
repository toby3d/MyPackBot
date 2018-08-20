package utils

import (
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

// SwitchButton helper just generate ReplyMarkup with SelfSwitch button
func SwitchButton(p *message.Printer) *tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitchSelf(
				p.Sprintf("select sticker"),
				" ",
			),
		),
	)
}
