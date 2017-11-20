package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

func commands(msg *telegram.Message) error {
	log.Ln("[commands] Check command message")
	switch strings.ToLower(msg.Command()) {
	case "start":
		log.Ln("[commands] Received a /start command")
		// TODO: Reply by greetings message and add user to DB
		return nil
	case "help":
		log.Ln("[commands] Received a /help command")
		// TODO: Reply by help instructions
		return nil
	case "addsticker":
		log.Ln("[commands] Received a /addsticker command")
		// TODO: Change current state to "addSticker" for adding sticker
		return nil
	case "delsticker":
		log.Ln("[commands] Received a /delsticker command")
		// TODO: Change current state to "delSticker" for deleting sticker
		return nil
	case "cancel":
		log.Ln("[commands] Received a /cancel command")
		// TODO: Change current state to default for aborting /addSticker or
		// /delSticker commands
		return nil
	default:
		return nil // Do nothing because unsupported command
	}
}
