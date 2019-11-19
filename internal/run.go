package internal

import (
	"time"

	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/handler"
	"gitlab.com/toby3d/mypackbot/internal/middleware"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (mpb *MyPackBot) Run() error {
	var updates tg.UpdatesChannel
	//nolint: godox
	/* TODO(toby3d)
	if mpb.config.IsSet("telegram.webhook") {
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
	}
	*/
	if mpb.config.IsSet("telegram.long_poll") {
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

	h := handler.NewHandler(mpb.store, mpb.usersStore, mpb.stickersStore).UpdateHandler
	chain := middleware.Chain{
		middleware.AcquireUser(mpb.usersStore),
		middleware.AcquireSticker(mpb.stickersStore),
		middleware.Hacktober(),
		middleware.Birthday(time.Date(0, time.November, 4, 0, 0, 0, 0, time.UTC)),
		middleware.UpdateLastSeen(mpb.usersStore),
	}
	h = chain.UpdateHandler(h)

	for update := range updates {
		update := update
		ctx := new(model.Context)
		ctx.Bot = mpb.bot
		ctx.Sticker = new(model.Sticker)
		ctx.Update = &update
		ctx.User = new(model.User)

		if err := h(ctx); err != nil {
			dlog.Ln("ERROR:", err.Error())
		}
	}

	return nil
}
