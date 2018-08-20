package commands

import (
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Help just send instructions about bot usage
func Help(msg *tg.Message) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	err := db.DB.ChangeUserState(msg.From.ID, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(
		msg.Chat.ID, p.Sprintf(
			"/%s - adds stickers one by one to your collection\n/%s - adds the entire set to your\n/%s at once - removes the sticker from your set one by one\n/%s - deletes the sticker set from your set\n/%s - removes all stickers from your set\n/%s - undoes the current operation\n\nTo view and send stickers from your set, simply type `@%s` (and space) in any chat.",
			models.CommandAddSticker,
			models.CommandAddPack,
			models.CommandDeleteSticker,
			models.CommandDeletePack,
			models.CommandReset,
			models.CommandCancel,
			bot.Bot.Username,
		),
	)
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.MenuKeyboard(p)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
