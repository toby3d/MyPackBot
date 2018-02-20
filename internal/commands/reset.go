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

// Reset prepare user to reset his pack
func Reset(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, total, err := db.UserStickers(msg.From.ID, 0, "")
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if total <= 0 {
		err = db.ChangeUserState(msg.From.ID, models.StateNone)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, T("error_already_reset"))
		reply.ParseMode = tg.ModeMarkdown
		reply.ReplyMarkup = helpers.MenuKeyboard(T)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	err = db.ChangeUserState(msg.From.ID, models.StateReset)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("reply_reset", map[string]interface{}{
		"KeyPhrase":     T("key_phrase"),
		"CancelCommand": models.CommandCancel,
	}))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.CancelButton(T)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
