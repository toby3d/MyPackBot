package handler

import (
	"context"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) isCallbackQuery(ctx context.Context, call *tg.CallbackQuery) (err error) {
	switch call.Data {
	case common.DataAddSticker:
		err = h.callbackAddSticker(ctx, call)
	case common.DataAddSet:
		err = h.callbackAddSet(ctx, call)
	case common.DataRemoveSticker:
		err = h.callbackRemoveSticker(ctx, call)
	case common.DataRemoveSet:
		err = h.callbackRemoveSet(ctx, call)
	}
	return err
}

func (h *Handler) callbackAddSticker(ctx context.Context, call *tg.CallbackQuery) (err error) {
	u, _ := ctx.Value("user").(*model.User)
	p, _ := ctx.Value("printer").(*message.Printer)
	s, _ := ctx.Value("sticker").(*model.Sticker)

	answer := tg.NewAnswerCallbackQuery(call.ID)
	answer.Text = p.Sprintf("callback__text_add-single", call.Message.ReplyToMessage.Sticker.SetName)
	if err = h.store.AddSticker(u, s); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-single"), common.DataRemoveSticker),
	))
	if len(call.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0], call.Message.ReplyMarkup.InlineKeyboard[0][1],
		)
	}

	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          call.Message.Chat.ID,
		InlineMessageID: call.InlineMessageID,
		MessageID:       call.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = h.bot.AnswerCallbackQuery(answer)
	return err
}

func (h *Handler) callbackAddSet(ctx context.Context, call *tg.CallbackQuery) (err error) {
	u, _ := ctx.Value("user").(*model.User)
	p, _ := ctx.Value("printer").(*message.Printer)
	s, _ := ctx.Value("sticker").(*model.Sticker)

	answer := tg.NewAnswerCallbackQuery(call.ID)
	set, err := h.bot.GetStickerSet(s.SetName)
	if err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)
		return err
	}

	for i := range set.Stickers {
		if s, err = h.store.Stickers().GetOrCreate(&model.Sticker{
			CreatedAt:  call.Message.Date,
			Emoji:      set.Stickers[i].Emoji,
			ID:         set.Stickers[i].FileID,
			IsAnimated: set.Stickers[i].IsAnimated,
			SetName:    set.Name,
		}); err != nil {
			answer.Text = "🐞 " + err.Error()
			_, err = h.bot.AnswerCallbackQuery(answer)
			return err
		}
		if err = h.store.AddSticker(u, s); err != nil {
			answer.Text = "🐞 " + err.Error()
			_, err = h.bot.AnswerCallbackQuery(answer)
			return err
		}
	}

	answer.Text = p.Sprintf("callback__text_add-set", set.Title)
	if err = h.store.AddStickersSet(u, set.Name); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-single"), common.DataRemoveSticker),
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-set"), common.DataRemoveSet),
	))
	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          call.Message.Chat.ID,
		InlineMessageID: call.InlineMessageID,
		MessageID:       call.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = h.bot.AnswerCallbackQuery(answer)
	return err
}

func (h *Handler) callbackRemoveSticker(ctx context.Context, call *tg.CallbackQuery) (err error) {
	u, _ := ctx.Value("user").(*model.User)
	p, _ := ctx.Value("printer").(*message.Printer)
	s, _ := ctx.Value("sticker").(*model.Sticker)

	answer := tg.NewAnswerCallbackQuery(call.ID)
	answer.Text = p.Sprintf("callback__text_remove-single")
	if err = h.store.RemoveSticker(u, s); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))
	if len(call.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0], call.Message.ReplyMarkup.InlineKeyboard[0][1],
		)
	}

	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          call.Message.Chat.ID,
		InlineMessageID: call.InlineMessageID,
		MessageID:       call.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = h.bot.AnswerCallbackQuery(answer)
	return err
}

func (h *Handler) callbackRemoveSet(ctx context.Context, call *tg.CallbackQuery) (err error) {
	u, _ := ctx.Value("user").(*model.User)
	p, _ := ctx.Value("printer").(*message.Printer)
	s, _ := ctx.Value("sticker").(*model.Sticker)

	answer := tg.NewAnswerCallbackQuery(call.ID)
	answer.Text = p.Sprintf("callback__text_remove-set", s.SetName)
	if err = h.store.RemoveStickersSet(u, s.SetName); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = h.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))
	if len(call.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-set"), common.DataAddSet),
		)
	}
	if _, err = h.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          call.Message.Chat.ID,
		InlineMessageID: call.InlineMessageID,
		MessageID:       call.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = h.bot.AnswerCallbackQuery(answer)
	return err
}
