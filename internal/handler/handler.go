package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/catalog"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/photos"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	up "gitlab.com/toby3d/mypackbot/internal/model/users/photos"
	us "gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
	"gitlab.com/toby3d/mypackbot/internal/store"
)

type Handler struct {
	users         users.Manager
	stickers      stickers.Manager
	photos        photos.Manager
	usersStickers us.Manager
	usersPhotos   up.Manager
	store         *store.Store
}

func NewHandler(us users.Manager, ss stickers.Manager, ps photos.Manager, uss us.Manager, ups up.Manager) *Handler {
	_ = catalog.RegisterPlurals()

	return &Handler{
		photos:        ps,
		stickers:      ss,
		users:         us,
		usersPhotos:   ups,
		usersStickers: uss,
		store: &store.Store{
			Photos:        ps,
			Stickers:      ss,
			UsersPhotos:   ups,
			UsersStickers: uss,
		},
	}
}

func (h *Handler) UpdateHandler(ctx *model.Context) (err error) {
	switch {
	case ctx.Request.IsMessage():
		err = h.IsMessage(ctx)
	case ctx.Request.IsCallbackQuery():
		err = h.IsCallbackQuery(ctx)
	case ctx.Request.IsInlineQuery():
		err = h.IsInlineQuery(ctx)
	}

	return err
}
