package main

import tg "github.com/toby3d/telegram" // My Telegram bindings

func commandCancel(msg *tg.Message) {
	bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	state, err := dbGetUserState(msg.From.ID)
	errCheck(err)

	var text string
	switch state {
	case stateAddSticker:
		text = T("cancel_add_sticker")
	case stateAddPack:
		text = T("cancel_add_pack")
	case stateDelete:
		text = T("cancel_del")
	case stateReset:
		text = T("cancel_reset")
	default:
		text = T("cancel_error")
	}

	err = dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	reply := tg.NewMessage(msg.Chat.ID, text)
	_, err = bot.SendMessage(reply)
	errCheck(err)
}
