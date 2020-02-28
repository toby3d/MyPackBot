package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
)

func AcquireUserSticker(store stickers.Reader) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) error {
		if ctx.Sticker == nil {
			return next(ctx)
		}

		ctx.HasSticker = store.Get(&model.UserSticker{
			UserID:    ctx.User.ID,
			StickerID: ctx.Sticker.ID,
		}) != nil

		return next(ctx)
	}
}
