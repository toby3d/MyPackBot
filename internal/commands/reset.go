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

// Reset prepare user to reset his pack
func Reset(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	stickers, err := db.DB.GetUserStickers(msg.From, &tg.InlineQuery{})
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if len(stickers) <= 0 {
		err = db.DB.ChangeUserState(msg.From, models.StateNone)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, T("error_already_reset"))
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyMarkup = utils.MenuKeyboard(T)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	err = db.DB.ChangeUserState(msg.From, models.StateReset)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("reply_reset", map[string]interface{}{
		"KeyPhrase":     T("key_phrase"),
		"CancelCommand": models.CommandCancel,
	}))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(T)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
