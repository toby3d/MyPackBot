package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) IsCallbackQuery(ctx *model.Context) (err error) {
	switch ctx.CallbackQuery.Data {
	case common.DataRemoveSticker, common.DataAddSticker:
		ctx.Set("is_add", ctx.CallbackQuery.Data == common.DataAddSticker)
		err = h.CallbackSticker(ctx)
	case common.DataAddSet, common.DataRemoveSet:
		ctx.Set("is_add", ctx.CallbackQuery.Data == common.DataAddSet)
		err = h.CallbackSet(ctx)
	}

	return err
}

func (h *Handler) CallbackSticker(ctx *model.Context) (err error) {
	isAdd, _ := ctx.Get("is_add").(bool)
	answer := tg.NewAnswerCallbackQuery(ctx.CallbackQuery.ID)
	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow())
	markup.InlineKeyboard[0] = make([]tg.InlineKeyboardButton, 1)

	if isAdd {
		answer.Text = ctx.T().Sprintf("callback__text_add-single")

		go h.store.AddSticker(ctx.User, ctx.Sticker)

		markup.InlineKeyboard[0][0] = tg.NewInlineKeyboardButton(
			ctx.T().Sprintf("sticker__button_remove-single"), common.DataRemoveSticker,
		)
	} else {
		answer.Text = ctx.T().Sprintf("callback__text_remove-single")

		go h.store.RemoveSticker(ctx.User, ctx.Sticker)

		markup.InlineKeyboard[0][0] = tg.NewInlineKeyboardButton(
			ctx.T().Sprintf("sticker__button_add-single"), common.DataAddSticker,
		)
	}

	if len(ctx.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0], ctx.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0][1],
		)
	}

	if _, err = ctx.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          ctx.CallbackQuery.Message.Chat.ID,
		InlineMessageID: ctx.CallbackQuery.InlineMessageID,
		MessageID:       ctx.CallbackQuery.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return ctx.Error(err)
	}

	_, err = ctx.AnswerCallbackQuery(answer)

	return ctx.Error(err)
}

func (h *Handler) CallbackSet(ctx *model.Context) (err error) {
	isAdd, _ := ctx.Get("is_add").(bool)
	answer := tg.NewAnswerCallbackQuery(ctx.CallbackQuery.ID)
	markup := new(tg.InlineKeyboardMarkup)

	if isAdd {
		set, err := ctx.GetStickerSet(ctx.Sticker.SetName)
		if err != nil {
			return ctx.Error(err)
		}

		answer.Text = ctx.T().Sprintf("callback__text_add-set", set.Title)

		go func(ctx *model.Context) {
			for i := range set.Stickers {
				sticker, err := h.stickersStore.GetOrCreate(&model.Sticker{
					CreatedAt:  ctx.CallbackQuery.Message.Date,
					UpdatedAt:  ctx.CallbackQuery.Message.Date,
					Width:      set.Stickers[i].Width,
					Height:     set.Stickers[i].Height,
					Emoji:      set.Stickers[i].Emoji,
					ID:         set.Stickers[i].FileID,
					IsAnimated: set.Stickers[i].IsAnimated,
					SetName:    set.Name,
				})
				if err != nil {
					continue
				}

				_ = h.store.AddSticker(ctx.User, sticker)
			}

			_ = h.store.AddStickersSet(ctx.User, set.Name)
		}(ctx)

		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				ctx.T().Sprintf("sticker__button_remove-single"), common.DataRemoveSticker,
			),
			tg.NewInlineKeyboardButton(
				ctx.T().Sprintf("sticker__button_remove-set"), common.DataRemoveSet,
			),
		))
	} else {
		answer.Text = ctx.T().Sprintf("callback__text_remove-set", ctx.Sticker.SetName)

		go h.store.RemoveStickersSet(ctx.User, ctx.Sticker.SetName)

		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				ctx.T().Sprintf("sticker__button_add-single"), common.DataAddSticker,
			),
			tg.NewInlineKeyboardButton(
				ctx.T().Sprintf("sticker__button_add-set"), common.DataAddSet,
			),
		))
	}

	if _, err = ctx.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          ctx.CallbackQuery.Message.Chat.ID,
		InlineMessageID: ctx.CallbackQuery.InlineMessageID,
		MessageID:       ctx.CallbackQuery.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return ctx.Error(err)
	}

	_, err = ctx.AnswerCallbackQuery(answer)

	return ctx.Error(err)
}
