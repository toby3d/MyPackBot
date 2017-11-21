package main

import (
	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

// message function check Message update on commands, sended stickers or other
// user stuff
func messages(msg *telegram.Message) {
	if msg.From.ID == bot.Self.ID {
		log.Ln("[messages] Received a message from myself, ignore this update")
		return
	}

	if msg.ForwardFrom != nil {
		if msg.ForwardFrom.ID == bot.Self.ID {
			log.Ln("[messages] Received a forward from myself, ignore this update")
			return
		}
	}

	if msg.IsCommand() {
		log.Ln("[message] Received a command message")
		commands(msg)
		return
	}

	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	if msg.Sticker != nil {
		state, err := dbGetUserState(msg.From.ID)
		errCheck(err)

		log.D(msg.Sticker)

		switch state {
		case stateNone:
			reply := telegram.NewMessage(
				msg.Chat.ID,              // chat
				T("error_command_first"), // text
			)
			_, err = bot.SendMessage(reply)
			errCheck(err)
		case stateAdding:
			exists, err := dbAddSticker(msg.From.ID, msg.Sticker.FileID, msg.Sticker.Emoji)
			errCheck(err)

			if exists {
				reply := telegram.NewMessage(
					msg.Chat.ID,     // chat
					T("add_exists"), // text
				)
				_, err = bot.SendMessage(reply)
				errCheck(err)

				_, err = dbChangeUserState(msg.From.ID, stateNone)
				errCheck(err)
				return
			}

			reply := telegram.NewMessage(
				msg.Chat.ID,      // chat
				T("add_success"), // text
			)
			_, err = bot.SendMessage(reply)
			errCheck(err)

			_, err = dbChangeUserState(msg.From.ID, stateNone)
			errCheck(err)
		case stateDeleting:
			notFound, err := dbDeleteSticker(msg.From.ID, msg.Sticker.FileID)
			errCheck(err)

			text := T("remove_success")
			if notFound {
				text = T("remove_already")
			}

			reply := telegram.NewMessage(
				msg.Chat.ID, // chat
				text,        // text
			)
			_, err = bot.SendMessage(reply)
			errCheck(err)

			_, err = dbChangeUserState(msg.From.ID, stateNone)
			errCheck(err)
		default:
			_, err := dbChangeUserState(msg.From.ID, stateNone)
			errCheck(err)

			messages(msg)
			return
		}
	}
}
