package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

const keyPhrase = "Yes, I am totally sure."

func commandCancel(msg *telegram.Message) {
	log.Ln("Received a /cancel command")
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	state, err := dbGetUserState(msg.From.ID)
	errCheck(err)

	text := T("cancel_error")
	switch state {
	case stateAddSticker:
		text = T("cancel_add_sticker")
	case stateAddPack:
		text = T("cancel_add_pack")
	case stateDelete:
		text = T("cancel_del")
	case stateReset:
		text = T("cancel_reset")
	}

	err = dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	reply := telegram.NewMessage(msg.Chat.ID, text)
	_, err = bot.SendMessage(reply)
	errCheck(err)
}
