package commands

import (
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Help just send instructions about bot usage
func Help(msg *tg.Message) {
	t, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	err = db.DB.ChangeUserState(msg.From, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(
		msg.Chat.ID, t("reply_help", map[string]interface{}{
			"AddStickerCommand":    models.CommandAddSticker,
			"AddPackCommand":       models.CommandAddPack,
			"DeleteStickerCommand": models.CommandDeleteSticker,
			"DeletePackCommand":    models.CommandDeletePack,
			"ResetCommand":         models.CommandReset,
			"CancelCommand":        models.CommandCancel,
			"Username":             bot.Bot.Username,
		}),
	)
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.MenuKeyboard(t)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
