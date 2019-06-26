package internal

import (
	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/events"
	"gitlab.com/toby3d/telegram"
)

func (mpb *MyPackBot) Run() error {
	switch {
	/* TODO
	case mpb.config.IsSet("telegram.webhook"):
		set := http.AcquireURI()
		defer http.ReleaseURI(set)

		cfg := mpb.config.Sub("telegram.webhook")
		mpb.updates = mpb.bot.NewWebhookChannel(
			set,
			&telegram.SetWebhookParameters{
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
		mpb.updates = mpb.bot.NewLongPollingChannel(&telegram.GetUpdatesParameters{
			AllowedUpdates: cfg.GetStringSlice("allowed_updates"),
			Limit:          cfg.GetInt("limit"),
			Offset:         cfg.GetInt("offset"),
			Timeout:        cfg.GetInt("timeout"),
		})
	}

	e := events.New(mpb.store)
	for update := range mpb.updates {
		var err error
		switch {
		case update.IsMessage():
			err = e.Message(mpb.bot, update.Message)
		case update.IsCallbackQuery():
			err = e.CallbackQuery(mpb.bot, update.CallbackQuery)
		case update.IsInlineQuery():
			err = e.InlineQuery(mpb.bot, update.InlineQuery)
		case update.IsChosenInlineResult():
			err = e.ChosenInlineResult(mpb.bot, update.ChosenInlineResult)
		default:
			dlog.D(update)
		}
		if err != nil {
			dlog.Ln("ERROR:", err.Error())
		}
	}
	return nil
}
