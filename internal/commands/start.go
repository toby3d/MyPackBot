package commands

import (
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/helpers"
	"github.com/toby3d/MyPackBot/internal/i18n"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// Start just send introduction about bot to user
func Start(msg *tg.Message) {
	err := db.DB.ChangeUserState(msg.From, models.StateNone)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	if msg.HasCommandArgument() {
		log.Ln("Received a", msg.Command(), "command with", msg.CommandArgument(), "argument")
		if strings.EqualFold(msg.CommandArgument(), models.CommandHelp) {
			Help(msg)
			return
		}
	}

	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	reply := tg.NewMessage(
		msg.Chat.ID,
		T("reply_start", map[string]interface{}{
			"Username": bot.Bot.Username,
			"ID":       bot.Bot.ID,
		}),
	)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.MenuKeyboard(T)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
