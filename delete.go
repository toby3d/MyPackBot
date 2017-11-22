package main

import (
	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func commandDelete(msg *telegram.Message) {
	log.Ln("Received a /remove command")
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	stickers, err := dbGetUserStickers(msg.From.ID, 0, "")
	errCheck(err)

	if len(stickers) <= 0 {
		reply := telegram.NewMessage(msg.Chat.ID, T("error_empty_remove"))
		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	err = dbChangeUserState(msg.From.ID, stateDelete)
	errCheck(err)

	markup := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonSwitchSelf(
				T("button_remove"),
				" ",
			),
		),
	)

	reply := telegram.NewMessage(msg.Chat.ID, T("reply_remove"))
	reply.ParseMode = telegram.ModeMarkdown
	reply.ReplyMarkup = &markup

	_, err = bot.SendMessage(reply)
	errCheck(err)
}

func actionDelete(msg *telegram.Message) {
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	notExist, err := dbDeleteSticker(msg.From.ID, msg.Sticker.FileID)
	errCheck(err)

	reply := telegram.NewMessage(msg.Chat.ID, T("success_remove"))
	reply.ParseMode = telegram.ModeMarkdown

	if notExist {
		reply.Text = T("error_already_remove")
	}

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
