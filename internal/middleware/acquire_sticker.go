package middleware

import (
	"context"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

func AcquireSticker(bot *tg.Bot, store store.StickersManager) Interceptor {
	return func(ctx context.Context, update *tg.Update, next model.UpdateFunc) (err error) {
		var s *model.Sticker
		switch {
		case update.IsMessage():
			if !update.Message.IsSticker() {
				return next(ctx, update)
			}

			s = utils.ConvertStickerToModel(update.Message.Sticker)
			s.CreatedAt = update.Message.Date
			s.UpdatedAt = update.Message.Date
		case update.IsCallbackQuery():
			if !update.CallbackQuery.Message.ReplyToMessage.IsSticker() ||
				!update.CallbackQuery.Message.IsReply() {
				return next(ctx, update)
			}

			s = utils.ConvertStickerToModel(update.CallbackQuery.Message.ReplyToMessage.Sticker)
			s.CreatedAt = update.CallbackQuery.Message.ReplyToMessage.Date
			s.UpdatedAt = update.CallbackQuery.Message.ReplyToMessage.Date
		default:
			return next(ctx, update)
		}

		if s.SetName != "" {
			go func() {
				set, err := bot.GetStickerSet(s.SetName)
				if err != nil {
					return
				}

				for _, setSticker := range set.Stickers {
					setSticker := setSticker
					sticker := utils.ConvertStickerToModel(&setSticker)
					sticker.CreatedAt = s.CreatedAt
					sticker.UpdatedAt = s.UpdatedAt
					_, _ = store.GetOrCreate(sticker)
				}
			}()
		}

		if s, err = store.GetOrCreate(s); err != nil {
			return err
		}

		return next(context.WithValue(ctx, common.ContextSticker, s), update)
	}
}
