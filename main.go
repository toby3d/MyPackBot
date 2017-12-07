package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

// bot is general structure of the bot
var bot *telegram.Bot

// main function is a general function for work of this bot
func main() {
	log.Ln("Let'g Get It Started...")
	var err error

	go dbInit()

	log.Ln("Initializing new bot via checking access_token...")
	bot, err = telegram.NewBot(cfg.UString("telegram.token"))
	errCheck(err)

	log.Ln("Let's check updates channel!")
	for update := range getUpdatesChannel() {
		switch {
		case update.InlineQuery != nil:
			// Just don't check same updates
			log.D(update.InlineQuery.Query)
			if len(update.InlineQuery.Query) > 25 {
				continue
			}

			inlineQuery(update.InlineQuery)
		case update.Message != nil:
			if update.Message.From.ID == bot.Self.ID {
				log.Ln("Received a message from myself, ignore this update")
				return
			}

			if update.Message.ForwardFrom != nil {
				if update.Message.ForwardFrom.ID == bot.Self.ID {
					log.Ln("Received a forward from myself, ignore this update")
					return
				}
			}

			messages(update.Message)
		default:
			log.Ln("Get unsupported update")
		}
		continue
	}

	err = db.Close()
	errCheck(err)
}
