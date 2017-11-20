package main

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

func commands(msg *telegram.Message) error {
	log.Ln("[commands] Check command message")
	var err error
	switch strings.ToLower(msg.Command()) {
	case "start":
		log.Ln("[commands] Received a /start command")
		// TODO: Reply by greetings message and add user to DB
		_, _, err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			fmt.Sprint("Hello, ", msg.From.FirstName, "!"), // text
		)
		_, err = bot.SendMessage(reply)
	case "help":
		log.Ln("[commands] Received a /help command")
		_, _, err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			fmt.Sprintln( // text
				"/start",
				"/help",
				"/addSticker",
				"/delSticker",
				"/cancel",
			),
		)
		_, err = bot.SendMessage(reply)
	case "addsticker":
		log.Ln("[commands] Received a /addsticker command")
		_, _, err = dbChangeUserState(msg.From.ID, stateAdding)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			"Send me any sticker for adding them in your pack.", // text
		)
		_, err = bot.SendMessage(reply)
	case "delsticker":
		log.Ln("[commands] Received a /delsticker command")
		_, _, err = dbChangeUserState(msg.From.ID, stateDeleting)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			"Send me sticker from your pack for remove them.", // text
		)
		_, err = bot.SendMessage(reply)
	case "cancel":
		log.Ln("[commands] Received a /cancel command")
		prev, _, err := dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		text := "What are you doing?!"
		switch prev {
		case stateAdding:
			prev = "You canceled adding a sticker to the set."
		case stateDeleting:
			prev = "You canceled the removal of the sticker from the set."
		case stateNone:
			prev = "Nothing to cancel."
		}

		reply := telegram.NewMessage(
			msg.Chat.ID, // chat
			text,        // text
		)
		_, err = bot.SendMessage(reply)
	default:
		log.Ln("[commands] Received unsupported command")
		// Do nothing because unsupported command
	}

	return err
}
