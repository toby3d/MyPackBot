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

// Cancel just cancel current user operation
func Cancel(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	state, err := db.UserState(msg.From.ID)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	var text string
	switch state {
	case models.StateAddSticker:
		text = T("cancel_add_sticker")
	case models.StateAddPack:
		text = T("cancel_add_pack")
	case models.StateDeleteSticker:
		text = T("cancel_del_sticker")
	case models.StateDeletePack:
		text = T("cancel_del_pack")
	case models.StateReset:
		text = T("cancel_reset")
	default:
		text = T("cancel_error")
	}

	err = db.ChangeUserState(msg.From.ID, models.StateNone)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = helpers.MenuKeyboard(T)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
