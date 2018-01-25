package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func commandStart(msg *tg.Message) {
	err := dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	if msg.HasCommandArgument() {
		log.Ln("Received a", msg.Command(), "command with", msg.CommandArgument(), "argument")
		if strings.ToLower(msg.CommandArgument()) == strings.ToLower(cmdHelp) {
			commandHelp(msg)
			return
		}
	}

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	reply := tg.NewMessage(msg.Chat.ID, T("reply_start", map[string]interface{}{
		"Username": bot.Self.Username,
		"ID":       bot.Self.ID,
	}))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = getMenuKeyboard(T)

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
