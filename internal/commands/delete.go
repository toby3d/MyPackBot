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

// Delete prepare user to remove some stickers or sets from his pack
func Delete(msg *tg.Message, pack bool) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, total, err := db.UserStickers(msg.From.ID, 0, "")
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if total <= 0 {
		err = db.ChangeUserState(msg.From.ID, models.StateNone)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, T("error_empty_del"))
		reply.ReplyMarkup = helpers.MenuKeyboard(T)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	reply := tg.NewMessage(msg.Chat.ID, T("reply_del_sticker"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.CancelButton(T)

	err = db.ChangeUserState(msg.From.ID, models.StateDeleteSticker)
	errors.Check(err)

	if pack {
		err = db.ChangeUserState(msg.From.ID, models.StateDeletePack)
		errors.Check(err)

		reply.Text = T("reply_del_pack")
	}

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply = tg.NewMessage(msg.Chat.ID, T("reply_switch_button"))
	reply.ReplyMarkup = helpers.SwitchButton(T)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
