package utils

import (
	"fmt"

	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

// CancelButton helper just generate ReplyMarkup with cancel button
func CancelButton(p *message.Printer) (rkm *tg.ReplyKeyboardMarkup) {
	rkm = tg.NewReplyKeyboardMarkup(
		tg.NewReplyKeyboardRow(
			tg.NewReplyKeyboardButton(
				fmt.Sprintf("❌ %s", p.Sprintf("cancel")),
			),
		),
	)
	rkm.ResizeKeyboard = true
	return
}
