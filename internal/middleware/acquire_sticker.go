package middleware

import (
	"context"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
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

			s = &model.Sticker{
				ID:         update.Message.Sticker.FileID,
				IsAnimated: update.Message.Sticker.IsAnimated,
				SetName:    update.Message.Sticker.SetName,
				Emoji:      update.Message.Sticker.Emoji,
				CreatedAt:  update.Message.Date,
			}
		case update.IsCallbackQuery():
			if !update.CallbackQuery.Message.IsReply() ||
				!update.CallbackQuery.Message.ReplyToMessage.IsSticker() {
				return next(ctx, update)
			}

			s = &model.Sticker{
				ID:         update.CallbackQuery.Message.ReplyToMessage.Sticker.FileID,
				IsAnimated: update.CallbackQuery.Message.ReplyToMessage.Sticker.IsAnimated,
				SetName:    update.CallbackQuery.Message.ReplyToMessage.Sticker.SetName,
				Emoji:      update.CallbackQuery.Message.ReplyToMessage.Sticker.Emoji,
				CreatedAt:  update.CallbackQuery.Message.ReplyToMessage.Date,
			}
		default:
			return next(ctx, update)
		}

		if s.SetName != "" {
			go func() {
				set, err := bot.GetStickerSet(s.SetName)
				if err != nil {
					return
				}

				for _, sticker := range set.Stickers {
					store.GetOrCreate(&model.Sticker{
						ID:         sticker.FileID,
						IsAnimated: sticker.IsAnimated,
						SetName:    set.Name,
						Emoji:      sticker.Emoji,
						CreatedAt:  s.CreatedAt,
					})
				}
			}()
		}

		s, err = store.GetOrCreate(s)
		if err != nil {
			return err
		}

		return next(context.WithValue(ctx, "sticker", s), update)
	}
}
