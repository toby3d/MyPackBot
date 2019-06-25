package internal

import (
	"github.com/kirillDanshin/dlog"
	http "github.com/valyala/fasthttp"
	"gitlab.com/toby3d/mypackbot/internal/update"
	"gitlab.com/toby3d/telegram"
)

func (mpb *MyPackBot) Run() error {
	switch {
	case mpb.config.IsSet("telegram.webhook"):
		cfg := mpb.config.Sub("telegram.webhook")

		set := http.AcquireURI()
		defer http.ReleaseURI(set)

		mpb.updates = mpb.bot.NewWebhookChannel(
			set,
			&telegram.SetWebhookParameters{
				AllowedUpdates: cfg.GetStringSlice("allowed_updates"),
			},
			cfg.GetString("certificate"),
			cfg.GetString("key"),
			cfg.GetString("serve"),
		)
	case mpb.config.IsSet("telegram.long_poll"):
		cfg := mpb.config.Sub("telegram.long_poll")

		if _, err := mpb.bot.DeleteWebhook(); err != nil {
			return err
		}

		mpb.updates = mpb.bot.NewLongPollingChannel(&telegram.GetUpdatesParameters{
			AllowedUpdates: cfg.GetStringSlice("allowed_updates"),
			Limit:          cfg.GetInt("limit"),
			Offset:         cfg.GetInt("offset"),
			Timeout:        cfg.GetInt("timeout"),
		})
	}

	for upd := range mpb.updates {

		var err error
		switch {
		case upd.IsMessage():
			err = update.Message(mpb.bot, mpb.userStore, mpb.stickerStore, upd.Message)
		case upd.IsCallbackQuery():
			err = update.CallbackQuery(mpb.bot, mpb.userStore, mpb.stickerStore, upd.CallbackQuery)
		case upd.IsInlineQuery():
			err = update.InlineQuery(mpb.bot, mpb.userStore, mpb.stickerStore, upd.InlineQuery)
		case upd.IsChosenInlineResult():
			err = update.ChosenInlineResult(mpb.bot, mpb.userStore, mpb.stickerStore, upd.ChosenInlineResult)
		default:
			dlog.D(upd)
		}
		if err != nil {
			dlog.Ln("ERROR:", err.Error())
		}
	}
	return nil
}
