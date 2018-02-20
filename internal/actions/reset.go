package actions

import (
	"strings"

	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/helpers"
	"github.com/toby3d/MyPackBot/internal/i18n"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// Reset action checks key phrase and reset user's pack
func Reset(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	err = db.ChangeUserState(msg.From.ID, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if !strings.EqualFold(msg.Text, T("key_phrase")) {
		reply := tg.NewMessage(msg.Chat.ID, T("error_reset_phrase"))
		reply.ParseMode = tg.ModeMarkdown
		reply.ReplyMarkup = helpers.MenuKeyboard(T)

		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	err = db.ResetUser(msg.From.ID)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_reset"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.MenuKeyboard(T)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
