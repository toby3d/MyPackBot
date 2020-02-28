package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/users/photos"
)

func AcquireUserPhoto(store photos.Reader) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) error {
		if ctx.Photo == nil {
			return next(ctx)
		}

		ctx.HasPhoto = store.Get(&model.UserPhoto{
			UserID:  ctx.User.ID,
			PhotoID: ctx.Photo.ID,
		}) != nil

		return next(ctx)
	}
}
