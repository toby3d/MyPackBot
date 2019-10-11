package internal

import (
	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/event"
	tg "gitlab.com/toby3d/telegram"
)

func (mpb *MyPackBot) Run() error {
	var updates tg.UpdatesChannel
	switch {
	/*
		case mpb.config.IsSet("telegram.webhook"):
			set := http.AcquireURI()
			defer http.ReleaseURI(set)

			cfg := mpb.config.Sub("telegram.webhook")
			updates = mpb.bot.NewWebhookChannel(
				set,
				&tg.SetWebhookParameters{
					AllowedUpdates: cfg.GetStringSlice("allowed_updates"),
				},
				cfg.GetString("certificate"),
				cfg.GetString("key"),
				cfg.GetString("serve"),
			)
	*/
	case mpb.config.IsSet("telegram.long_poll"):
		if _, err := mpb.bot.DeleteWebhook(); err != nil {
			return err
		}

		cfg := mpb.config.Sub("telegram.long_poll")
		updates = mpb.bot.NewLongPollingChannel(&tg.GetUpdatesParameters{
			AllowedUpdates: cfg.GetStringSlice("allowed_updates"),
			Limit:          cfg.GetInt("limit"),
			Offset:         cfg.GetInt("offset"),
			Timeout:        cfg.GetInt("timeout"),
		})
	}

	e := event.NewEvent(mpb.bot, mpb.store)
	for update := range updates {
		var err error
		switch {
		case update.IsMessage():
			err = e.Message(update.Message)
		case update.IsCallbackQuery():
			err = e.CallbackQuery(update.CallbackQuery)
		case update.IsInlineQuery():
			err = e.InlineQuery(update.InlineQuery)
		// case update.IsChosenInlineResult():
		// 	err = e.ChosenInlineResult(update.ChosenInlineResult)
		default:
			dlog.D(update)
		}
		if err != nil {
			dlog.Ln("ERROR:", err.Error())
		}
	}
	return nil
}
