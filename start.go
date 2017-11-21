package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func commandStart(msg *telegram.Message) {
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	err := dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	log.Ln("Received a /start command")
	if msg.HasArgument() {
		if strings.ToLower(msg.CommandArgument()) == "add" {
			commandAdd(msg)
			return
		}
	}

	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	reply := telegram.NewMessage(
		msg.Chat.ID, T("reply_start", map[string]interface{}{
			"Username": bot.Self.Username,
		}),
	)
	reply.ParseMode = telegram.ModeMarkdown

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
