package updates

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/config"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// Channel return webhook or long polling channel with bot updates
func Channel(webhookMode bool) tg.UpdatesChannel {
	log.Ln("Preparing channel for updates...")
	if !webhookMode {
		log.Ln("Use LongPolling updates")

		info, err := bot.Bot.GetWebhookInfo()
		errors.Check(err)

		if info.URL != "" {
			log.Ln("Deleting webhook...")
			_, err := bot.Bot.DeleteWebhook()
			errors.Check(err)
		}

		return bot.Bot.NewLongPollingChannel(&tg.GetUpdatesParameters{
			Offset:         0,
			Limit:          100,
			Timeout:        60,
			AllowedUpdates: models.AllowedUpdates,
		})
	}

	set := config.Config.GetString("telegram.webhook.set")
	listen := config.Config.GetString("telegram.webhook.listen")
	serve := config.Config.GetString("telegram.webhook.serve")

	log.Ln(
		"Trying set webhook on address:",
		fmt.Sprint(set, listen, bot.Bot.AccessToken),
	)

	log.Ln("Creating new webhook...")
	webhook := tg.NewWebhook(fmt.Sprint(set, listen, bot.Bot.AccessToken), nil)
	webhook.MaxConnections = 40
	webhook.AllowedUpdates = models.AllowedUpdates

	return bot.Bot.NewWebhookChannel(
		webhook, // params
		"",      // certFile
		"",      // keyFile
		set,     // set
		fmt.Sprint(listen, bot.Bot.AccessToken), // listen
		serve, // serve
	)
}
