package commands

import (
	"strings"

	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Start just send introduction about bot to user
func Start(msg *tg.Message) {
	err := db.DB.ChangeUserState(msg.From.ID, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if msg.HasCommandArgument() {
		log.Ln("Received a", msg.Command(), "command with", msg.CommandArgument(), "argument")
		if strings.EqualFold(msg.CommandArgument(), tg.CommandHelp) {
			Help(msg)
			return
		}
	}

	p := utils.NewPrinter(msg.From.LanguageCode)

	reply := tg.NewMessage(
		msg.Chat.ID,
		p.Sprintf(
			"Hello, I'm @%s!\nI can create your personal set of stickers from other sets.\nWithout restrictions and installation. In any chat. Is free.",
			bot.Bot.Username,
		),
	)
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.MenuKeyboard(p)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
