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

	chain := middleware.Chain{
		middleware.AcquireUser(mpb.users),
		middleware.ChatAction(),
		middleware.AcquirePhoto(mpb.photos),
		middleware.AcquireUserPhoto(mpb.usersPhotos),
		middleware.AcquireSticker(mpb.stickers),
		middleware.AcquireUserSticker(mpb.usersStickers),
		middleware.Birthday(time.Date(0, time.November, 4, 0, 0, 0, 0, time.UTC)),
		middleware.Hacktober(),
		middleware.UpdateLastSeen(mpb.users),
	}
	h := chain.UpdateHandler(handler.NewHandler(
		mpb.users,
		mpb.stickers,
		mpb.photos,
		mpb.usersStickers,
		mpb.usersPhotos,
	).UpdateHandler)

	for update := range updates {
		update := update
		go func(update *tg.Update) {
			ctx := new(model.Context)
			ctx.Bot = mpb.bot
			ctx.Request = update
			if err := h(ctx); err != nil {
				dlog.D(err)
			}
		}(&update)
	}

	return nil
}
