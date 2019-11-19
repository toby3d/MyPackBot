package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

func AcquireSticker(bot *tg.Bot, store store.StickersManager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		ctx.Sticker = new(model.Sticker)
		switch {
		case ctx.IsMessage():
			if !ctx.Message.IsSticker() {
				return next(ctx)
			}

			ctx.Sticker = utils.ConvertStickerToModel(ctx.Message.Sticker)
			ctx.Sticker.CreatedAt = ctx.Message.Date
			ctx.Sticker.UpdatedAt = ctx.Message.Date
		case ctx.IsCallbackQuery():
			if !ctx.CallbackQuery.Message.ReplyToMessage.IsSticker() ||
				!ctx.CallbackQuery.Message.IsReply() {
				return next(ctx)
			}

			ctx.Sticker = utils.ConvertStickerToModel(ctx.CallbackQuery.Message.ReplyToMessage.Sticker)
			ctx.Sticker.CreatedAt = ctx.CallbackQuery.Message.ReplyToMessage.Date
			ctx.Sticker.UpdatedAt = ctx.CallbackQuery.Message.ReplyToMessage.Date
		default:
			return next(ctx)
		}

		if ctx.Sticker.SetName != "" {
			go func() {
				set, err := bot.GetStickerSet(ctx.Sticker.SetName)
				if err != nil {
					return
				}

				for _, setSticker := range set.Stickers {
					setSticker := setSticker
					sticker := utils.ConvertStickerToModel(&setSticker)
					sticker.CreatedAt = ctx.Sticker.CreatedAt
					sticker.UpdatedAt = ctx.Sticker.UpdatedAt
					_, _ = store.GetOrCreate(sticker)
				}
			}()
		}

		if ctx.Sticker, err = store.GetOrCreate(ctx.Sticker); err != nil {
			return err
		}

		return next(ctx)
	}
}
