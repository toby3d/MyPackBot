package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

// bot is general structure of the bot
var bot *telegram.Bot

// main function is a general function for work of this bot
func main() {
	log.Ln("[main] Let'g Get It Started...")
	var err error
	log.Ln("[main] Initializing new bot via checking access_token...")
	bot, err = telegram.NewBot(bot.AccessToken)
	errCheck(err)

	log.Ln("[main] Initializing channel for updates...")
	updates, err := getUpdatesChannel()
	errCheck(err)

	log.Ln("[main] Let's check updates channel!")
	for update := range updates {
		switch {
		case update.ChosenInlineResult != nil:
			log.Ln("[main] Get ChosenInlineResult update")
			// TODO: Save info in Yandex.AppMetrika
		case update.InlineQuery != nil:
			log.Ln("[main] Get InlineQuery update")
			// TODO: Search stickers via inline
		case update.Message != nil:
			log.Ln("[main] Get Message update")
			// TODO: Added support of commands, grab and save sticker in DB
		default:
			log.Ln("[main] Get unsupported update")
			continue // Ignore any other updates
		}
	}
}
