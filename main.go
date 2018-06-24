package main

import (
	"flag"

	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	"gitlab.com/toby3d/mypackbot/internal/updates"
)

var flagWebhook = flag.Bool(
	"webhook", false,
	"enable work via webhooks (required valid certificates)",
)

// init prepare configuration and other things for successful start
func init() {
	log.Ln("Initializing...")

	// Preload localization strings
	err := i18n.Open("i18n/")
	errors.Check(err)

	// Preload configuration file
	config.Open("configs/config.yaml")

	// Open database or create new one
	db.Open("stickers.db")

	// Create bot with credentials from config
	bot.Bot, err = bot.New(config.Config.GetString("telegram.token"))
	errors.Check(err)
}

// main function is a general function for work of this bot
func main() {
	flag.Parse() // Parse flagWebhook

	channel, err := updates.Channel(*flagWebhook)
	errors.Check(err)

	for update := range channel {
		log.D(update)
		switch {
		case update.IsInlineQuery():
			updates.InlineQuery(update.InlineQuery)
		case update.IsMessage():
			updates.Message(update.Message)
		case update.IsChannelPost():
			updates.ChannelPost(update.ChannelPost)
		}
	}
}
