package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
)

func AcquireUserSticker(store stickers.Manager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) error {
		if ctx.Sticker == nil {
			return next(ctx)
		}

		ctx.UserSticker = store.Get(&model.UserSticker{
			UserID:    ctx.User.ID,
			StickerID: ctx.Sticker.ID,
		})

		return next(ctx)
	}
}
