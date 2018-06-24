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

// Cancel just cancel current user operation
func Cancel(msg *tg.Message) {
	t, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	state, err := db.DB.GetUserState(msg.From)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	var text string
	switch state {
	case models.StateAddSticker:
		text = t("cancel_add_sticker")
	case models.StateAddPack:
		text = t("cancel_add_pack")
	case models.StateDeleteSticker:
		text = t("cancel_del_sticker")
	case models.StateDeletePack:
		text = t("cancel_del_pack")
	case models.StateReset:
		text = t("cancel_reset")
	default:
		text = t("cancel_error")
	}

	err = db.DB.ChangeUserState(msg.From, models.StateNone)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = utils.MenuKeyboard(t)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
