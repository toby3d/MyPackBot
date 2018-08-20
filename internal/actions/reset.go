package actions

import (
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Reset action checks key phrase and reset user's pack
func Reset(msg *tg.Message) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	err := db.DB.ChangeUserState(msg.From.ID, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if !strings.EqualFold(msg.Text, p.Sprintf("Yes, I am absolutely sure.")) {
		reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("Wrong phrase for reset. The action was canceled."))
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyMarkup = utils.MenuKeyboard(p)

		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	err = db.DB.ResetUser(msg.From.ID)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("Your set has successfully reseted!"))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.MenuKeyboard(p)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
