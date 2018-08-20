package actions

import (
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Error action send error reply about invalid user request
func Error(msg *tg.Message) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	_, err := bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(
		msg.Chat.ID, p.Sprintf(
			"I have no idea what to do with this sticker.\nPlease use /%s, /%s, /%s or /%s first.",
			models.CommandAddSticker,
			models.CommandAddPack,
			models.CommandDeleteSticker,
			models.CommandDeletePack,
		),
	)
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.MenuKeyboard(p)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
