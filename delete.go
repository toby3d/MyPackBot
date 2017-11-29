package main

import "github.com/toby3d/go-telegram" // My Telegram bindings

func commandDelete(msg *telegram.Message) {
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	stickers, err := dbGetUserStickers(msg.From.ID, "")
	errCheck(err)

	if len(stickers) <= 0 {
		err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

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

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

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
