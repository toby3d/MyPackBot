package main

import (
	"fmt"
	"time"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

const keyPhrase = "Yes, I am totally sure."

func commandReset(msg *telegram.Message) {
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	stickers, err := dbGetUserStickers(msg.From.ID, 0, "")
	errCheck(err)

	if len(stickers) <= 0 {
		reply := telegram.NewMessage(msg.Chat.ID, T("error_already_reset"))
		reply.ParseMode = telegram.ModeMarkdown

		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	err = dbChangeUserState(msg.From.ID, stateReset)
	errCheck(err)

	reply := telegram.NewMessage(
		msg.Chat.ID,
		T("reply_reset", map[string]interface{}{
			"KeyPhrase":     keyPhrase,
			"CancelCommand": cmdCancel,
		}))
	reply.ParseMode = telegram.ModeMarkdown

	_, err = bot.SendMessage(reply)
	errCheck(err)
}

func actionReset(msg *telegram.Message) {
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	err := dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	if msg.Text != keyPhrase {
		reply := telegram.NewMessage(msg.Chat.ID, T("error_reset_phrase"))
		reply.ParseMode = telegram.ModeMarkdown

		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	err = dbResetUserStickers(msg.From.ID)
	errCheck(err)

	reply := telegram.NewMessage(msg.Chat.ID, T("success_reset"))
	reply.ParseMode = telegram.ModeMarkdown

	_, err = bot.SendMessage(reply)
	errCheck(err)

	for i := 1; i <= 3; i++ {
		bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

		text := T(fmt.Sprint("meta_reset_", i))

		time.Sleep(time.Minute * time.Duration(len(text)) / 1000)

		reply = telegram.NewMessage(msg.Chat.ID, text)
		reply.ParseMode = telegram.ModeMarkdown

		_, err = bot.SendMessage(reply)
		errCheck(err)
	}
}
