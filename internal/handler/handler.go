package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/photos"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	up "gitlab.com/toby3d/mypackbot/internal/model/users/photos"
	us "gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
	"gitlab.com/toby3d/mypackbot/internal/store"
)

type Handler struct {
	users         users.ReadWriter
	stickers      stickers.ReadWriter
	photos        photos.ReadWriter
	usersStickers us.ReadWriter
	usersPhotos   up.ReadWriter
	store         *store.Store
}

func NewHandler(us users.ReadWriter, ss stickers.ReadWriter, ps photos.ReadWriter, uss us.ReadWriter,
	ups up.ReadWriter) *Handler {
	return &Handler{
		photos:        ps,
		stickers:      ss,
		users:         us,
		usersPhotos:   ups,
		usersStickers: uss,
		store:         store.NewStore(uss, ups),
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
