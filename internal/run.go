package internal

import (
	"net"
	"time"

	"github.com/kirillDanshin/dlog"
	http "github.com/valyala/fasthttp"
	"gitlab.com/toby3d/mypackbot/internal/handler"
	"gitlab.com/toby3d/mypackbot/internal/middleware"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (mpb *MyPackBot) Run() error {
	updates, shutdown, err := mpb.getUpdateChannel()
	if err != nil {
		return err
	}

	defer func() { _ = shutdown() }()

	chain := middleware.Chain{
		middleware.AcquireUser(mpb.users),
		middleware.AcquirePrinter(),
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

func (mpb *MyPackBot) getUpdateChannel() (tg.UpdatesChannel, tg.ShutdownFunc, error) {
	switch {
	case mpb.config.IsSet("telegram.webhook"):
		cfg := mpb.config.Sub("telegram.webhook")
		u := http.AcquireURI()

		defer http.ReleaseURI(u)

		cert := make([]string, 0, 2)
		if cfg.IsSet("certificate") {
			cert = append(cert, cfg.GetString("certificate"))
		}

		if cfg.IsSet("key") {
			cert = append(cert, cfg.GetString("key"))
		}

		ln, err := net.Listen("tcp", cfg.GetString("serve"))
		if err != nil {
			return nil, nil, err
		}

		updates, shutdown := mpb.bot.NewWebhookChannel(u, &tg.SetWebhookParameters{
			AllowedUpdates: cfg.GetStringSlice("allowed_updates"),
		}, ln, cert...)

		return updates, shutdown, nil
	case mpb.config.IsSet("telegram.long_poll"):
		if _, err := mpb.bot.DeleteWebhook(); err != nil {
			return nil, nil, err
		}

		cfg := mpb.config.Sub("telegram.long_poll")
		updates := mpb.bot.NewLongPollingChannel(&tg.GetUpdatesParameters{
			AllowedUpdates: cfg.GetStringSlice("allowed_updates"),
			Limit:          cfg.GetInt("limit"),
			Offset:         cfg.GetInt("offset"),
			Timeout:        cfg.GetInt("timeout"),
		})

		return updates, tg.ShutdownFunc(nil), nil
	}

	return nil, nil, nil
}
