package main

import tg "github.com/toby3d/telegram" // My Telegram bindings

func commandHelp(msg *tg.Message) {
	bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

	err := dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	reply := tg.NewMessage(
		msg.Chat.ID, T("reply_help", map[string]interface{}{
			"AddStickerCommand":    cmdAddSticker,
			"AddPackCommand":       cmdAddPack,
			"DeleteStickerCommand": cmdDeleteSticker,
			"DeletePackCommand":    cmdDeletePack,
			"ResetCommand":         cmdReset,
			"CancelCommand":        cmdCancel,
			"Username":             bot.Self.Username,
		}),
	)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = getMenuKeyboard(T)

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
