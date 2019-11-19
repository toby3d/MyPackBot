package handler

import (
	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
)

type Handler struct {
	store         store.Manager
	usersStore    users.Manager
	stickersStore stickers.Manager
}

func NewHandler(store store.Manager, usersStore users.Manager, stickersStore stickers.Manager) *Handler {
	return &Handler{
		store:         store,
		usersStore:    usersStore,
		stickersStore: stickersStore,
	}
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

	return ctx.Error(err)
}
