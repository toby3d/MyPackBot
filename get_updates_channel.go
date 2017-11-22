package main

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog"     // Insert logs only in debug builds
	json "github.com/pquerna/ffjson/ffjson" // Fastest JSON unmarshalling
	"github.com/toby3d/go-telegram"         // My Telegram bindings
	http "github.com/valyala/fasthttp"      // Fastest http-requests
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

	tgHookSet := cfg.UString("telegram.webhook.set")
	tgHookListen := cfg.UString("telegram.webhook.listen")
	tgHookServe := cfg.UString("telegram.webhook.serve")

	log.Ln(
		"Trying set webhook on address:",
		fmt.Sprint(tgHookSet, tgHookListen, bot.AccessToken),
	)

	log.Ln("Creating new webhook...")
	webhook := telegram.NewWebhook(
		fmt.Sprint(tgHookSet, tgHookListen, bot.AccessToken), nil,
	)
	webhook.AllowedUpdates = allowedUpdates

	log.Ln("Setting new webhook...")
	_, err = bot.SetWebhook(webhook)
	errCheck(err)

	channel := make(chan telegram.Update, 100)
	go func() {
		log.Ln("Listen and serve...")
		err := http.ListenAndServe(
			tgHookServe,
			func(ctx *http.RequestCtx) {
				log.Ln("Catch request on path:", string(ctx.Path()))
				if !strings.HasPrefix(
					string(ctx.Path()), fmt.Sprint(tgHookListen, bot.AccessToken),
				) {
					return
				}

				log.Ln("Catch supported request:")
				log.Ln(string(ctx.Request.Body()))

				if ctx.Request.Body() == nil {
					return
				}

				var update telegram.Update
				err := json.Unmarshal(ctx.Request.Body(), &update)
				errCheck(err)

				log.Ln("Unmarshled next:")
				log.D(update)

				channel <- update
			},
		)
		errCheck(err)
	}()

	return channel
}
