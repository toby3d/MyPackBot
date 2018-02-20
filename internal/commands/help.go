package commands

import (
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/helpers"
	"github.com/toby3d/MyPackBot/internal/i18n"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// Help just send instructions about bot usage
func Help(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	err = db.ChangeUserState(msg.From.ID, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(
		msg.Chat.ID, T("reply_help", map[string]interface{}{
			"AddStickerCommand":    models.CommandAddSticker,
			"AddPackCommand":       models.CommandAddPack,
			"DeleteStickerCommand": models.CommandDeleteSticker,
			"DeletePackCommand":    models.CommandDeletePack,
			"ResetCommand":         models.CommandReset,
			"CancelCommand":        models.CommandCancel,
			"Username":             bot.Bot.Self.Username,
		}),
	)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.MenuKeyboard(T)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
