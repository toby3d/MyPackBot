package actions

import (
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Reset action checks key phrase and reset user's pack
func Reset(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	err = db.DB.ChangeUserState(msg.From, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if !strings.EqualFold(msg.Text, T("key_phrase")) {
		reply := tg.NewMessage(msg.Chat.ID, T("error_reset_phrase"))
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyMarkup = utils.MenuKeyboard(T)

		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	err = db.DB.ResetUser(msg.From)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_reset"))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.MenuKeyboard(T)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
