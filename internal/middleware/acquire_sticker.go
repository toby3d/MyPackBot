package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	tg "gitlab.com/toby3d/telegram"
)

func AcquireSticker(store stickers.Manager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		ctx.Sticker = new(model.Sticker)
		switch {
		case ctx.IsMessage() && ctx.Message.IsSticker():
			ctx.Sticker = stickerToModel(ctx.Message.Sticker)
			ctx.Sticker.CreatedAt = ctx.Message.Date
			ctx.Sticker.UpdatedAt = ctx.Message.Date
		case ctx.IsMessage() && ctx.Message.IsReply() && ctx.Message.ReplyToMessage.IsSticker():
			ctx.Sticker = stickerToModel(ctx.Message.ReplyToMessage.Sticker)
			ctx.Sticker.CreatedAt = ctx.Message.Date
			ctx.Sticker.UpdatedAt = ctx.Message.Date
		case ctx.IsCallbackQuery():
			if !ctx.CallbackQuery.Message.IsReply() ||
				!ctx.CallbackQuery.Message.ReplyToMessage.IsSticker() {
				return next(ctx)
			}

			ctx.Sticker = stickerToModel(ctx.CallbackQuery.Message.ReplyToMessage.Sticker)
			ctx.Sticker.CreatedAt = ctx.CallbackQuery.Message.ReplyToMessage.Date
			ctx.Sticker.UpdatedAt = ctx.CallbackQuery.Message.ReplyToMessage.Date
		default:
			return next(ctx)
		}

		if ctx.Sticker.InSet() {
			go func() {
				set, err := ctx.GetStickerSet(ctx.Sticker.SetName)
				if err != nil {
					return
				}

				for _, setSticker := range set.Stickers {
					setSticker := setSticker
					sticker := stickerToModel(&setSticker)
					sticker.CreatedAt = ctx.Sticker.CreatedAt
					sticker.UpdatedAt = ctx.Sticker.UpdatedAt
					_ = store.Create(sticker)
				}
			}()
		}

		if ctx.Sticker, err = store.GetOrCreate(ctx.Sticker); err != nil {
			return err
		}

		return next(ctx)
	}
}

func stickerToModel(s *tg.Sticker) *model.Sticker {
	sticker := new(model.Sticker)
	sticker.ID = s.FileID
	sticker.Emoji = s.Emoji
	sticker.Width = s.Width
	sticker.Height = s.Height
	sticker.IsAnimated = s.IsAnimated
	sticker.SetName = s.SetName

	if !sticker.InSet() {
		sticker.SetName = common.SetNameUploaded
	}

	return sticker
}
