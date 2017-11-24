package main

import (
	"fmt"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

// allowedUpdates is a value for parameter of updates configuration
var allowedUpdates = []string{
	telegram.UpdateChosenInlineResult, // For collecting statistics
	telegram.UpdateInlineQuery,        // For searching and sending stickers
	telegram.UpdateMessage,            // For get commands and messages
}

// getUpdatesChannel return webhook or long polling channel with bot updates
func getUpdatesChannel() telegram.UpdatesChannel {
	log.Ln("Preparing channel for updates...")

	log.Ln("Deleting webhook if exists")
	_, err := bot.DeleteWebhook()
	errCheck(err)

	if !*flagWebhook {
		log.Ln("Use LongPolling updates")
		return bot.NewLongPollingChannel(&telegram.GetUpdatesParameters{
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
	webhook := telegram.NewWebhook(
		fmt.Sprint(set, listen, bot.AccessToken), nil,
	)
	webhook.MaxConnections = 100
	webhook.AllowedUpdates = allowedUpdates

	return bot.NewWebhookChannel(webhook, "", "", set, listen, serve)
}
