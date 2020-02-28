package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/photos"
	tg "gitlab.com/toby3d/telegram"
)

func AcquirePhoto(store photos.ReadWriter) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		switch {
		case ctx.Request.IsMessage():
			switch {
			case ctx.Request.Message.IsPhoto():
				ctx.Photo = photoToModel(ctx.Request.Message.Photo)
				ctx.Photo.CreatedAt = ctx.Request.Message.Date
				ctx.Photo.UpdatedAt = ctx.Request.Message.Date
			case ctx.Request.Message.IsReply() && ctx.Request.Message.ReplyToMessage.IsPhoto():
				ctx.Photo = photoToModel(ctx.Request.Message.ReplyToMessage.Photo)
				ctx.Photo.CreatedAt = ctx.Request.Message.Date
				ctx.Photo.UpdatedAt = ctx.Request.Message.Date
			default:
				return next(ctx)
			}
		case ctx.Request.IsCallbackQuery():
			if !ctx.Request.CallbackQuery.Message.IsReply() ||
				!ctx.Request.CallbackQuery.Message.ReplyToMessage.IsPhoto() {
				return next(ctx)
			}

			ctx.Photo = photoToModel(ctx.Request.CallbackQuery.Message.ReplyToMessage.Photo)
			ctx.Photo.CreatedAt = ctx.Request.CallbackQuery.Message.ReplyToMessage.Date
			ctx.Photo.UpdatedAt = ctx.Request.CallbackQuery.Message.ReplyToMessage.Date
		default:
			return next(ctx)
		}

		if ctx.Photo, err = store.GetOrCreate(ctx.Photo); err != nil {
			return err
		}

		return next(ctx)
	}
}

func photoToModel(photoSize tg.Photo) *model.Photo {
	p := photoSize[len(photoSize)-1]
	photo := new(model.Photo)
	photo.ID = p.FileUniqueID
	photo.Width = p.Width
	photo.Height = p.Height
	photo.FileID = p.FileID

	return photo
}
