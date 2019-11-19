package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) IsCallbackQuery(ctx *model.Context) (err error) {
	switch ctx.CallbackQuery.Data {
	case common.DataAddSticker:
		err = h.CallbackAddSticker(ctx)
	case common.DataAddSet:
		err = h.CallbackAddSet(ctx)
	case common.DataRemoveSticker:
		err = h.CallbackRemoveSticker(ctx)
	case common.DataRemoveSet:
		err = h.CallbackRemoveSet(ctx)
	}

	return err
}

func (h *Handler) CallbackAddSticker(ctx *model.Context) (err error) {
	answer := tg.NewAnswerCallbackQuery(ctx.CallbackQuery.ID)
	answer.Text = ctx.Printer.Sprintf("callback__text_add-single")

	if err = h.store.AddSticker(ctx.User, ctx.Sticker); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)

		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(ctx.Printer.Sprintf("sticker__button_remove-single"), common.DataRemoveSticker),
	))

	if len(ctx.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0], ctx.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0][1],
		)
	}

	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          ctx.CallbackQuery.Message.Chat.ID,
		InlineMessageID: ctx.CallbackQuery.InlineMessageID,
		MessageID:       ctx.CallbackQuery.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return err
	}

	_, err = h.bot.AnswerCallbackQuery(answer)

	return err
}

func (h *Handler) CallbackAddSet(ctx *model.Context) (err error) {
	answer := tg.NewAnswerCallbackQuery(ctx.CallbackQuery.ID)

	set, err := h.bot.GetStickerSet(ctx.Sticker.SetName)
	if err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)

		return err
	}

	answer.Text = ctx.Printer.Sprintf("callback__text_add-set", set.Title)

	for i := range set.Stickers {
		if ctx.Sticker, err = h.store.Stickers().GetOrCreate(&model.Sticker{
			CreatedAt:  ctx.CallbackQuery.Message.Date,
			UpdatedAt:  ctx.CallbackQuery.Message.Date,
			Width:      set.Stickers[i].Width,
			Height:     set.Stickers[i].Height,
			Emoji:      set.Stickers[i].Emoji,
			ID:         set.Stickers[i].FileID,
			IsAnimated: set.Stickers[i].IsAnimated,
			SetName:    set.Name,
		}); err != nil {
			answer.Text = "🐞 " + err.Error()
			_, err = h.bot.AnswerCallbackQuery(answer)

			return err
		}

		_ = h.store.AddSticker(ctx.User, ctx.Sticker)
	}

	if err = h.store.AddStickersSet(ctx.User, set.Name); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)

		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(ctx.Printer.Sprintf("sticker__button_remove-single"), common.DataRemoveSticker),
		tg.NewInlineKeyboardButton(ctx.Printer.Sprintf("sticker__button_remove-set"), common.DataRemoveSet),
	))

	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          ctx.CallbackQuery.Message.Chat.ID,
		InlineMessageID: ctx.CallbackQuery.InlineMessageID,
		MessageID:       ctx.CallbackQuery.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return err
	}

	_, err = h.bot.AnswerCallbackQuery(answer)

	return err
}

func (h *Handler) CallbackRemoveSticker(ctx *model.Context) (err error) {
	answer := tg.NewAnswerCallbackQuery(ctx.CallbackQuery.ID)
	answer.Text = ctx.Printer.Sprintf("callback__text_remove-single")

	if err = h.store.RemoveSticker(ctx.User, ctx.Sticker); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)

		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(ctx.Printer.Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))

	if len(ctx.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0], ctx.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0][1],
		)
	}

	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          ctx.CallbackQuery.Message.Chat.ID,
		InlineMessageID: ctx.CallbackQuery.InlineMessageID,
		MessageID:       ctx.CallbackQuery.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return err
	}

	_, err = h.bot.AnswerCallbackQuery(answer)

	return err
}

func (h *Handler) CallbackRemoveSet(ctx *model.Context) (err error) {
	answer := tg.NewAnswerCallbackQuery(ctx.CallbackQuery.ID)
	answer.Text = ctx.Printer.Sprintf("callback__text_remove-set", ctx.Sticker.SetName)

	if err = h.store.RemoveStickersSet(ctx.User, ctx.Sticker.SetName); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)

		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(ctx.Printer.Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))

	if len(ctx.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(ctx.Printer.Sprintf("sticker__button_add-set"), common.DataAddSet),
		)
	}

	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          ctx.CallbackQuery.Message.Chat.ID,
		InlineMessageID: ctx.CallbackQuery.InlineMessageID,
		MessageID:       ctx.CallbackQuery.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return err
	}

	_, err = h.bot.AnswerCallbackQuery(answer)

	return err
}
