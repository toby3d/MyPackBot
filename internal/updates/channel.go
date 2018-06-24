package updates

import (
	"fmt"
	"net/url"

	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
)

// Channel return webhook or long polling channel with bot updates
func Channel(webhookMode bool) (updates tg.UpdatesChannel, err error) {
	log.Ln("Preparing channel for updates...")
	if !webhookMode {
		log.Ln("Use LongPolling updates")

		var info *tg.WebhookInfo
		info, err = bot.Bot.GetWebhookInfo()
		if err != nil {
			return
		}

		if info.URL != "" {
			log.Ln("Deleting webhook...")
			_, err = bot.Bot.DeleteWebhook()
			return
		}

		updates = bot.Bot.NewLongPollingChannel(&tg.GetUpdatesParameters{
			Offset:         0,
			Limit:          100,
			Timeout:        60,
			AllowedUpdates: models.AllowedUpdates,
		})
		return
	}

	set, err := url.Parse(config.Config.GetString("telegram.webhook.set"))
	if err != nil {
		return nil, err
	}

	listen := config.Config.GetString("telegram.webhook.listen")
	serve := config.Config.GetString("telegram.webhook.serve")

	log.Ln(
		"Trying set webhook on address:",
		fmt.Sprint(set.String(), bot.Bot.AccessToken),
	)

	log.Ln("Creating new webhook...")
	params := tg.NewWebhook(fmt.Sprint(set, listen, bot.Bot.AccessToken), nil)
	params.MaxConnections = 40
	params.AllowedUpdates = models.AllowedUpdates

	updates = bot.Bot.NewWebhookChannel(
		set,
		params, // params
		"",     // certFile
		"",     // keyFile
		serve,  // serve
	)

	return
}
