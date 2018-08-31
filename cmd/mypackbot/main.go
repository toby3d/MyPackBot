package main

//go:generate gotext -srclang=en update -out=i18n_gen.go -lang=en

import (
	"flag"
	"path/filepath"

	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/updates"
)

var (
	flagWebhook = flag.Bool(
		"webhook", false,
		"enable work via webhooks (required valid certificates)",
	)
	flagConfig = flag.String(
		"config",
		filepath.Join(".", "config", "config.yaml"),
		"set specific path to config",
	)
	flagDB = flag.String(
		"db",
		filepath.Join(".", "stickers.db"),
		"set specific path to stickers database",
	)
)

// init prepare configuration and other things for successful start
func init() {
	log.Ln("Initializing...")
	var err error

	// Preload configuration file
	config.Config, err = config.Open(*flagConfig)
	errors.Check(err)

	// Open database or create new one
	db.DB, err = db.Open(*flagDB)
	errors.Check(err)

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
