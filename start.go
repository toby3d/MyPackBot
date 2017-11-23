package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

func commandStart(msg *telegram.Message) {
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	err := dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	if msg.HasArgument() {
		if strings.ToLower(msg.CommandArgument()) == cmdAddSticker {
			commandAdd(msg, false)
			return
		}
	}

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	reply := telegram.NewMessage(
		msg.Chat.ID, T("reply_start", map[string]interface{}{
			"Username": bot.Self.Username,
		}),
	)
	reply.ParseMode = telegram.ModeMarkdown

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
