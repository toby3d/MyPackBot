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

func NewHandler(b *tg.Bot, s store.Manager) *Handler {
	return &Handler{
		bot:   b,
		store: s,
	}
}

func (h *Handler) UpdateHandler(ctx context.Context, upd *tg.Update) (err error) {
	switch {
	case upd.IsMessage():
		err = h.isMessage(ctx, upd.Message)
	case upd.IsCallbackQuery():
		err = h.isCallbackQuery(ctx, upd.CallbackQuery)
	case upd.IsInlineQuery():
		err = h.isInlineQuery(ctx, upd.InlineQuery)
	default:
		dlog.D(upd)
	}
	return err
}
