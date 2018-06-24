package actions

import (
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Error action send error reply about invalid user request
func Error(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(
		msg.Chat.ID, T("error_unknown", map[string]interface{}{
			"AddStickerCommand":    models.CommandAddSticker,
			"AddPackCommand":       models.CommandAddPack,
			"DeleteStickerCommand": models.CommandDeleteSticker,
			"DeletePackCommand":    models.CommandDeletePack,
		}),
	)
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.MenuKeyboard(T)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
