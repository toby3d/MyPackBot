package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

// message function check Message update on commands, sended stickers or other
// user stuff
func messages(msg *telegram.Message) error {
	if msg.From.ID == bot.Self.ID ||
		msg.ForwardFrom.ID == bot.Self.ID {
		log.Ln("[messages] Received a message from myself, ignore this update")
		return nil
	}

	if msg.IsCommand() {
		log.Ln("[message] Received a command message")
		return commands(msg)
	}

	if msg.Sticker != nil {
		state, err := dbGetUserState(msg.From.ID)
		errCheck(err)

		switch state {
		case stateNone:
			reply := telegram.NewMessage(
				msg.Chat.ID, // chat
				"Use /addSticker or /delSticker command first.", // text
			)
			_, err = bot.SendMessage(reply)
		case stateAdding:
			err = dbAddSticker(msg.From.ID, msg.Sticker.FileID, msg.Sticker.Emoji)
			errCheck(err)

			reply := telegram.NewMessage(
				msg.Chat.ID, // chat
				"Added!",    // text
			)
			_, err = bot.SendMessage(reply)
		case stateDeleting:
			err = dbDeleteSticker(msg.From.ID, msg.Sticker.FileID)
			errCheck(err)

			reply := telegram.NewMessage(
				msg.Chat.ID, // chat
				"Removed!",  // text
			)
			_, err = bot.SendMessage(reply)
		default:
			_, _, err := dbChangeUserState(msg.From.ID, stateNone)
			errCheck(err)
			return messages(msg)
		}
		// TODO: Upload new or delete exist sticker in pack
		return err
	}

	return nil // Do nothing because unsupported actions
}
