package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) IsCallbackQuery(ctx *model.Context) (err error) {
	if !ctx.Request.IsCallbackQuery() {
		return nil
	}

	switch ctx.Request.CallbackQuery.Data {
	case common.DataAdd, common.DataAddSet:
		err = h.CallbackAdd(ctx)
	case common.DataDel, common.DataDelSet:
		err = h.CallbackDel(ctx)
	}

	if err != nil {
		return err
	}

	_, err = ctx.AnswerCallbackQuery(&tg.AnswerCallbackQueryParameters{
		CallbackQueryID: ctx.Request.CallbackQuery.ID,
	})

	return err
}

func (h *Handler) CallbackAdd(ctx *model.Context) (err error) {
	if !ctx.Request.IsCallbackQuery() {
		return err
	}

	editMessage := new(tg.EditMessageReplyMarkupParameters)
	editMessage.ChatID = ctx.Request.CallbackQuery.Message.Chat.ID
	editMessage.InlineMessageID = ctx.Request.CallbackQuery.InlineMessageID
	editMessage.MessageID = ctx.Request.CallbackQuery.Message.ID

	switch {
	case ctx.Photo != nil:
		if err = h.CommandAddPhoto(ctx); err != nil {
			return err
		}

		editMessage.ReplyMarkup = h.GetPhotoKeyboard(ctx)
	case ctx.Sticker != nil:
		if ctx.Request.CallbackQuery.Data == common.DataAddSet {
			err = h.CommandAddSet(ctx)
		} else {
			err = h.CommandAddSticker(ctx)
		}

		if err != nil {
			return err
		}

		editMessage.ReplyMarkup = h.GetStickerKeyboard(ctx)
	default:
		return err
	}

	_, err = ctx.EditMessageReplyMarkup(editMessage)

	return err
}

func (h *Handler) CallbackDel(ctx *model.Context) (err error) {
	if !ctx.Request.IsCallbackQuery() {
		return err
	}

	editMessage := new(tg.EditMessageReplyMarkupParameters)
	editMessage.ChatID = ctx.Request.CallbackQuery.Message.Chat.ID
	editMessage.InlineMessageID = ctx.Request.CallbackQuery.InlineMessageID
	editMessage.MessageID = ctx.Request.CallbackQuery.Message.ID

	switch {
	case ctx.Photo != nil:
		if err = h.CommandDelPhoto(ctx); err != nil {
			return err
		}

		editMessage.ReplyMarkup = h.GetPhotoKeyboard(ctx)
	case ctx.Sticker != nil:
		if ctx.Request.CallbackQuery.Data == common.DataDelSet {
			err = h.CommandDelSet(ctx)
		} else {
			err = h.CommandDelSticker(ctx)
		}

		if err != nil {
			return err
		}

		editMessage.ReplyMarkup = h.GetStickerKeyboard(ctx)
	default:
		return err
	}

	_, err = ctx.EditMessageReplyMarkup(editMessage)

	return err
}
