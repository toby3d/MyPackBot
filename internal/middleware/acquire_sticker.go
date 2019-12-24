package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	tg "gitlab.com/toby3d/telegram"
)

func AcquireSticker(store stickers.Manager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		switch {
		case ctx.Request.IsMessage():
			switch {
			case ctx.Request.Message.IsSticker():
				ctx.Sticker = stickerToModel(ctx.Request.Message.Sticker)
				ctx.Sticker.CreatedAt = ctx.Request.Message.Date
				ctx.Sticker.UpdatedAt = ctx.Request.Message.Date
			case ctx.Request.Message.IsReply() && ctx.Request.Message.ReplyToMessage.IsSticker():
				ctx.Sticker = stickerToModel(ctx.Request.Message.ReplyToMessage.Sticker)
				ctx.Sticker.CreatedAt = ctx.Request.Message.Date
				ctx.Sticker.UpdatedAt = ctx.Request.Message.Date
			default:
				return next(ctx)
			}
		case ctx.Request.IsCallbackQuery():
			if !ctx.Request.CallbackQuery.Message.IsReply() ||
				!ctx.Request.CallbackQuery.Message.ReplyToMessage.IsSticker() {
				return next(ctx)
			}

			ctx.Sticker = stickerToModel(ctx.Request.CallbackQuery.Message.ReplyToMessage.Sticker)
			ctx.Sticker.CreatedAt = ctx.Request.CallbackQuery.Message.ReplyToMessage.Date
			ctx.Sticker.UpdatedAt = ctx.Request.CallbackQuery.Message.ReplyToMessage.Date
		default:
			return next(ctx)
		}

		if ctx.Sticker.InSet() {
			migrateSet(ctx, store)
		}

		if ctx.Sticker, err = store.GetOrCreate(ctx.Sticker); err != nil {
			return err
		}

		return next(ctx)
	}
}

func stickerToModel(s *tg.Sticker) *model.Sticker {
	sticker := new(model.Sticker)
	sticker.FileID = s.FileID
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

func migrateSet(ctx *model.Context, store stickers.Manager) {
	tgSet, err := ctx.GetStickerSet(ctx.Sticker.SetName)
	if err != nil || tgSet == nil || len(tgSet.Stickers) == 0 {
		stickers, _ := store.GetSet(ctx.Sticker.SetName)
		ctx.Sticker.SetName = common.SetNameUploaded

		go func() {
			for i := range stickers {
				stickers[i].SetName = ctx.Sticker.SetName
				stickers[i].UpdatedAt = ctx.Sticker.UpdatedAt
				_ = store.Update(stickers[i])
			}
		}()
	} else {
		ctx.Set("set_name", tgSet.Title)

		dbSet, _ := store.GetSet(tgSet.Name)
		for i := range tgSet.Stickers {
			for j := range dbSet {
				if tgSet.Stickers[i].FileID == dbSet[j].FileID {
					continue
				}

				_ = store.Create(stickerToModel(&tgSet.Stickers[i]))
			}
		}
	}
}
