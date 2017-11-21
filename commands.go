package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func commands(msg *telegram.Message) {
	log.Ln("[commands] Check command message")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	switch strings.ToLower(msg.Command()) {
	case "start":
		log.Ln("[commands] Received a /start command")
		// TODO: Reply by greetings message and add user to DB
		_, err := dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			T("start_message", map[string]interface{}{
				"Username": bot.Self.Username,
			}), // text
		)
		_, err = bot.SendMessage(reply)
		errCheck(err)
	case "help":
		log.Ln("[commands] Received a /help command")
		_, err := dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			T("help_message", map[string]interface{}{
				"Username": bot.Self.Username,
			}), // text
		)
		_, err = bot.SendMessage(reply)
		errCheck(err)
	case "add":
		log.Ln("[commands] Received a /add command")
		_, err := dbChangeUserState(msg.From.ID, stateAdding)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID,    // chat
			T("add_reply"), // text
		)
		_, err = bot.SendMessage(reply)
		errCheck(err)
	case "remove":
		log.Ln("[commands] Received a /remove command")
		_, err := dbChangeUserState(msg.From.ID, stateDeleting)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID,       // chat
			T("remove_reply"), // text
		)
		_, err = bot.SendMessage(reply)
		errCheck(err)
	case "cancel":
		log.Ln("[commands] Received a /cancel command")
		prev, err := dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		text := T("error_cancel_nothing")
		switch prev {
		case stateAdding:
			prev = T("add_cancel")
		case stateDeleting:
			prev = T("remove_cancel")
		}

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			text,        // text
		)
		_, err = bot.SendMessage(reply)
		errCheck(err)
	default:
		log.Ln("[commands] Received unsupported command")
		// Do nothing because unsupported command
	}
}
