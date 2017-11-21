package main

import (
	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func commandHelp(msg *telegram.Message) {
	log.Ln("Received a /help command")
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	err := dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	markup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonSwitch(
				T("button_share"),
				" ",
			),
		),
	)

	reply := telegram.NewMessage(
		msg.Chat.ID, T("reply_help", map[string]interface{}{
			"Username": bot.Self.Username,
		}),
	)
	reply.ParseMode = telegram.ModeMarkdown
	reply.ReplyMarkup = &markup

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
