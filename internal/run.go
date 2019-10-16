package internal

import (
	"context"
	"time"

	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/handler"
	"gitlab.com/toby3d/mypackbot/internal/middleware"
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

	h := handler.NewHandler(mpb.bot, mpb.store)
	bDay := time.Date(0, time.November, 4, 0, 0, 0, 0, time.UTC)
	chain := middleware.Chain{
		middleware.AcquireUser(mpb.store.Users()),
		middleware.AcquirePrinter(),
		middleware.AcquireSticker(mpb.bot, mpb.store.Stickers()),
		middleware.Birthday(mpb.bot, mpb.store.Users(), bDay),
		middleware.UpdateLastSeen(mpb.store.Users()),
	}
	for update := range updates {
		if err := chain.UpdateHandler(h.UpdateHandler)(context.Background(), &update); err != nil {
			dlog.Ln("ERROR:", err.Error())
		}
	}
	return nil
}
