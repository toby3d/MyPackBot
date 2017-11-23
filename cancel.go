package main

import "github.com/toby3d/go-telegram" // My Telegram bindings

func commandCancel(msg *telegram.Message) {
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

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

	reply := telegram.NewMessage(msg.Chat.ID, text)
	_, err = bot.SendMessage(reply)
	errCheck(err)
}
