package main

import tg "github.com/toby3d/telegram" // My Telegram bindings

func commandDelete(msg *tg.Message) {
	bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	_, total, err := dbGetUserStickers(msg.From.ID, 0, "")
	errCheck(err)

	if total <= 0 {
		err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		reply := tg.NewMessage(msg.Chat.ID, T("error_empty_del"))
		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	err = dbChangeUserState(msg.From.ID, stateDelete)
	errCheck(err)

	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitchSelf(
				T("button_del"),
				" ",
			),
		),
	)

	reply := tg.NewMessage(msg.Chat.ID, T("reply_del"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = &markup

	_, err = bot.SendMessage(reply)
	errCheck(err)
}

func actionDelete(msg *tg.Message) {
	bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	notExist, err := dbDeleteSticker(
		msg.From.ID,
		msg.Sticker.SetName,
		msg.Sticker.FileID,
	)
	errCheck(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_del"))
	reply.ParseMode = tg.ModeMarkdown

	if notExist {
		reply.Text = T("error_already_del")
	}

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
