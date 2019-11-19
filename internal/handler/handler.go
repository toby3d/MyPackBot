package handler

import (
	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
)

type Handler struct {
	store store.Manager
}

func NewHandler(store store.Manager) *Handler {
	return &Handler{store: store}
}

func (h *Handler) UpdateHandler(ctx *model.Context) (err error) {
	switch {
	case ctx.IsMessage():
		err = h.IsMessage(ctx)
	case ctx.IsCallbackQuery():
		err = h.IsCallbackQuery(ctx)
	case ctx.IsInlineQuery():
		err = h.IsInlineQuery(ctx)
	default:
		dlog.D(ctx)
	}

	return err
}
