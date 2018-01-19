package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram"     // My Telegram bindings
)

// bot is general structure of the bot
var bot *tg.Bot

// main function is a general function for work of this bot
func main() {
	log.Ln("Let'g Get It Started...")
	var err error

	go dbInit()

	log.Ln("Initializing new bot via checking access_token...")
	bot, err = tg.NewBot(cfg.UString("telegram.token"))
	errCheck(err)

	log.Ln("Let's check updates channel!")
	for update := range getUpdatesChannel() {
		switch {
		case update.InlineQuery != nil:
			log.D(update.InlineQuery)
			updateInlineQuery(update.InlineQuery)
		case update.Message != nil:
			log.D(update.Message)
			updateMessage(update.Message)
		case update.ChannelPost != nil:
			log.D(update.ChannelPost)
			updateChannelPost(update.ChannelPost)
		default:
			log.D(update)
		}
	}

	err = db.Close()
	errCheck(err)
}
