package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/users/photos"
)

func AcquireUserPhoto(store photos.Manager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) error {
		if ctx.Photo == nil {
			return next(ctx)
		}

		ctx.UserPhoto = store.Get(&model.UserPhoto{
			UserID:  ctx.User.ID,
			PhotoID: ctx.Photo.ID,
		})

		return next(ctx)
	}
}
