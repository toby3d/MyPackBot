package handler

import (
	"context"

	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
)

type Handler struct {
	bot   *tg.Bot
	store store.Manager
}

func NewHandler(bot *tg.Bot, store store.Manager) *Handler {
	return &Handler{
		bot:   bot,
		store: store,
	}
}

func (h *Handler) UpdateHandler(ctx context.Context, upd *tg.Update) (err error) {
	switch {
	case upd.IsMessage():
		err = h.IsMessage(ctx, upd.Message)
	case upd.IsCallbackQuery():
		err = h.IsCallbackQuery(ctx, upd.CallbackQuery)
	case upd.IsInlineQuery():
		err = h.IsInlineQuery(ctx, upd.InlineQuery)
	default:
		dlog.D(upd)
	}

	return err
}
