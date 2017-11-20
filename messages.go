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
		// TODO: Upload new or delete exist sticker in pack
		return nil
	}

	return nil // Do nothing because unsupported actions
}
