package main

import (
	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func commandAdd(msg *telegram.Message) {
	log.Ln("Received a /add command")
	log.Ln("Change", msg.From.ID, "state to", stateAdd)
	err := dbChangeUserState(msg.From.ID, stateAdd)
	errCheck(err)

	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	reply := telegram.NewMessage(msg.Chat.ID, T("reply_add"))
	reply.ParseMode = telegram.ModeMarkdown

	log.Ln("Sending add reply...")
	_, err = bot.SendMessage(reply)
	errCheck(err)
}

func actionAdd(msg *telegram.Message) {
	log.Ln("Received a /add action")
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	log.Ln("Change", msg.From.ID, "state to", stateNone)
	err := dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	exists, err := dbAddSticker(msg.From.ID, msg.Sticker.FileID, msg.Sticker.Emoji)
	errCheck(err)

	markup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonSwitch(
				T("button_share"),
				msg.Sticker.Emoji,
			),
		),
	)

	reply := telegram.NewMessage(msg.Chat.ID, T("success_add"))
	reply.ParseMode = telegram.ModeMarkdown
	reply.ReplyMarkup = &markup

	if exists {
		reply.Text = T("error_already_add")
	}

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
