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

// Delete prepare user to remove some stickers or sets from his pack
func Delete(msg *tg.Message, pack bool) {
	t, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	stickers, err := db.DB.GetUserStickers(msg.From.ID, &tg.InlineQuery{})
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if len(stickers) <= 0 {
		err = db.DB.ChangeUserState(msg.From.ID, models.StateNone)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, t("error_empty_del"))
		reply.ReplyMarkup = utils.MenuKeyboard(t)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	reply := tg.NewMessage(msg.Chat.ID, t("reply_del_sticker"))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(t)

	err = db.DB.ChangeUserState(msg.From.ID, models.StateDeleteSticker)
	errors.Check(err)

	if pack {
		err = db.DB.ChangeUserState(msg.From.ID, models.StateDeletePack)
		errors.Check(err)

		reply.Text = t("reply_del_pack")
	}

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply = tg.NewMessage(msg.Chat.ID, t("reply_switch_button"))
	reply.ReplyMarkup = utils.SwitchButton(t)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
