package utils

import (
	"fmt"

	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

// MenuKeyboard helper just generate ReplyMarkup with menu buttons
func MenuKeyboard(p *message.Printer) (rkm *tg.ReplyKeyboardMarkup) {
	rkm = tg.NewReplyKeyboardMarkup(
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(
				fmt.Sprintf("➕ %s", p.Sprintf("add a sticker")),
			), tg.NewReplyKeyboardButton(
				fmt.Sprintf("📦 %s", p.Sprintf("add set")),
			),
		),
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(
				fmt.Sprintf("🗑 %s", p.Sprintf("remove sticker")),
			), tg.NewReplyKeyboardButton(
				fmt.Sprintf("🗑 %s", p.Sprintf("delete set")),
			),
		),
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(
				fmt.Sprintf("🔥 %s", p.Sprintf("reset set")),
			),
		),
	)

	rkm.ResizeKeyboard = true
	return
}
