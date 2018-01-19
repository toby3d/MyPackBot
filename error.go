package main

import tg "github.com/toby3d/telegram"

func actionError(msg *tg.Message) {
	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	reply := tg.NewMessage(
		msg.Chat.ID, T("error_unknown", map[string]interface{}{
			"AddStickerCommand":    cmdAddSticker,
			"AddPackCommand":       cmdAddPack,
			"DeleteStickerCommand": cmdDeleteSticker,
			"DeletePackCommand":    cmdDeletePack,
		}),
	)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = getMenuKeyboard(T)

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
