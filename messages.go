package main

import (
	"strings"

	tg "github.com/toby3d/telegram"
)

func messages(msg *tg.Message) {
	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	switch {
	case strings.EqualFold(msg.Text, T("button_add_sticker")):
		commandAdd(msg, false)
	case strings.EqualFold(msg.Text, T("button_add_pack")):
		commandAdd(msg, true)
	case strings.EqualFold(msg.Text, T("button_del_sticker")):
		commandDelete(msg, false)
	case strings.EqualFold(msg.Text, T("button_del_pack")):
		commandDelete(msg, true)
	case strings.EqualFold(msg.Text, T("button_reset")):
		commandReset(msg)
	case strings.EqualFold(msg.Text, T("button_cancel")):
		commandCancel(msg)
	case strings.EqualFold(msg.Text, T("meta_key_phrase")):
		actions(msg)
	}
}
