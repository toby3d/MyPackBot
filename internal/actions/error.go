package actions

import (
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/helpers"
	"github.com/toby3d/MyPackBot/internal/i18n"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
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
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.MenuKeyboard(T)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
