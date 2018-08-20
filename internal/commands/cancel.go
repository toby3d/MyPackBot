package commands

import (
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Cancel just cancel current user operation
func Cancel(msg *tg.Message) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	state, err := db.DB.GetUserState(msg.From.ID)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	var text string
	switch state {
	case models.StateAddSticker:
		text = p.Sprintf("You canceled the process of adding new stickers to your collection.")
	case models.StateAddPack:
		text = p.Sprintf("You canceled the process of adding new sets to yours.")
	case models.StateDeleteSticker:
		text = p.Sprintf("You canceled the process of removing the sticker from your collection.")
	case models.StateDeletePack:
		text = p.Sprintf("You canceled the process of removing sets from your collection.")
	case models.StateReset:
		text = p.Sprintf("You canceled the process of resetting your collection.")
	default:
		text = p.Sprintf("Nothing to cancel.")
	}

	err = db.DB.ChangeUserState(msg.From.ID, models.StateNone)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = utils.MenuKeyboard(p)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
