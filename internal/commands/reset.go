package commands

import (
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Reset prepare user to reset his pack
func Reset(msg *tg.Message) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	stickers, err := db.DB.GetUserStickers(msg.From.ID, &tg.InlineQuery{})
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if len(stickers) <= 0 {
		err = db.DB.ChangeUserState(msg.From.ID, models.StateNone)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("There is nothing to discard, the set is already empty."))
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyMarkup = utils.MenuKeyboard(p)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	err = db.DB.ChangeUserState(msg.From.ID, models.StateReset)
	errors.Check(err)

	keyPhrase := p.Sprintf("Yes, I am absolutely sure.")
	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf(
		"This operation will remove *all* stickers from your set and *this can not be undone*.\n\nWrite `%s` to confirm my intention to zero my brain (oh god why).\nIf you use /%s to cancel the current operation.",
		keyPhrase,
		models.CommandCancel,
	))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(p)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
