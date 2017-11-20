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
func getUpdatesChannel() (telegram.UpdatesChannel, error) {
	log.Ln("[getUpdatesChannel] Preparing channel for updates...")

	log.Ln("[getUpdatesChannel] Deleting webhook if exists")
	_, err := bot.DeleteWebhook()
	errCheck(err)

	if !*flagWebhook {
		log.Ln("[getUpdatesChannel] Use LongPolling updates")
		return bot.NewLongPollingChannel(&telegram.GetUpdatesParameters{
			Offset:         0,
			Limit:          100,
			Timeout:        60,
			AllowedUpdates: allowedUpdates,
		}), nil
	}

	log.Ln(
		"[getUpdatesChannel] Trying set webhook on address:",
		fmt.Sprint(tgHookSet, tgHookListen, bot.AccessToken),
	)

	log.Ln("[getUpdatesChannel] Creating new webhook...")
	webhook := telegram.NewWebhook(
		fmt.Sprint(tgHookSet, tgHookListen, bot.AccessToken),
		"cert.pem",
	)
	webhook.AllowedUpdates = allowedUpdates

	log.Ln("[getUpdatesChannel] Setting new webhook...")
	_, err = bot.SetWebhook(webhook)
	if err != nil {
		return nil, err
	}

	channel := make(chan telegram.Update, 100)
	go func() {
		log.Ln("[getUpdatesChannel] Listen and serve TLS...")
		err := http.ListenAndServeTLS(
			tgHookServe,
			"cert.pem",
			"cert.key",
			func(ctx *http.RequestCtx) {
				log.Ln(
					"[getUpdatesChannel] Catch request on path:",
					string(ctx.Path()),
				)
				if strings.HasPrefix(
					string(ctx.Path()),
					fmt.Sprint(tgHookListen, bot.AccessToken),
				) {
					log.Ln(
						"[getUpdatesChannel] Catch supported request:",
						string(ctx.Request.Body()),
					)

					var update telegram.Update
					err := json.Unmarshal(ctx.Request.Body(), &update)
					errCheck(err)

					log.Ln("[getUpdatesChannel] Unmarshled next:")
					log.D(update)

					channel <- update
				}
			},
		)
		errCheck(err)
	}()

	return channel, nil
}
