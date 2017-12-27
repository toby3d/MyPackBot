package main

import (
	"fmt"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram"     // My Telegram bindings
)

// allowedUpdates is a value for parameter of updates configuration
var allowedUpdates = []string{
	telegram.UpdateInlineQuery, // For searching and sending stickers
	telegram.UpdateMessage,     // For get commands and messages
}

// getUpdatesChannel return webhook or long polling channel with bot updates
func getUpdatesChannel() tg.UpdatesChannel {
	log.Ln("Preparing channel for updates...")
	if !*flagWebhook {
		log.Ln("Use LongPolling updates")

		log.Ln("Deleting webhook if exists")
		_, err := bot.DeleteWebhook()
		errCheck(err)

		return bot.NewLongPollingChannel(&tg.GetUpdatesParameters{
			Offset:         0,
			Limit:          100,
			Timeout:        60,
			AllowedUpdates: allowedUpdates,
		})
	}

	set := cfg.UString("telegram.webhook.set")
	listen := cfg.UString("telegram.webhook.listen")
	serve := cfg.UString("telegram.webhook.serve")

	log.Ln(
		"Trying set webhook on address:",
		fmt.Sprint(set, listen, bot.AccessToken),
	)

	log.Ln("Creating new webhook...")
	webhook := tg.NewWebhook(
		fmt.Sprint(set, listen, bot.AccessToken), nil,
	)
	webhook.MaxConnections = 40
	webhook.AllowedUpdates = allowedUpdates

	return bot.NewWebhookChannel(
		webhook, // params
		"",      // certFile
		"",      // keyFile
		set,     // set
		fmt.Sprint(listen, bot.AccessToken), // listen
		serve, // serve
	)
}
