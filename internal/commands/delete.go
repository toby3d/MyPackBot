package commands

import (
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Delete prepare user to remove some stickers or sets from his pack
func Delete(msg *tg.Message, pack bool) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	stickers, err := db.DB.GetUserStickers(msg.From.ID, &tg.InlineQuery{})
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if len(stickers) <= 0 {
		err = db.DB.ChangeUserState(msg.From.ID, models.StateNone)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("Nothing to delete, the set is already empty."))
		reply.ReplyMarkup = utils.MenuKeyboard(p)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("Send a sticker from your set to remove it."))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(p)

	err = db.DB.ChangeUserState(msg.From.ID, models.StateDeleteSticker)
	errors.Check(err)

	if pack {
		err = db.DB.ChangeUserState(msg.From.ID, models.StateDeletePack)
		errors.Check(err)

		reply.Text = p.Sprintf("Send a sticker from your set to remove all of its set.")
	}

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply = tg.NewMessage(msg.Chat.ID, p.Sprintf("This button will help you quickly call your kit to select the sticker you want."))
	reply.ReplyMarkup = utils.SwitchButton(p)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
