package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) IsMessage(ctx *model.Context) (err error) {
	switch {
	case ctx.Message.IsCommand():
		err = h.IsCommand(ctx)
	case ctx.Message.IsSticker():
		err = h.IsSticker(ctx)
	}

	return ctx.Error(err)
}

// IsSticker send to the Sticker a Message with CallbackQuery buttons for selecting actions.
func (h *Handler) IsSticker(ctx *model.Context) error {
	us, err := h.store.GetSticker(ctx.User, ctx.Sticker)
	if err != nil {
		return ctx.Error(err)
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(ctx.T().Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))

	if ctx.Sticker.InSet() {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(ctx.T().Sprintf("sticker__button_add-set"), common.DataAddSet),
		)
	}

	if us != nil {
		markup.InlineKeyboard[0][0] = tg.NewInlineKeyboardButton(
			ctx.T().Sprintf("sticker__button_remove-single"), common.DataRemoveSticker,
		)

		if ctx.Sticker.InSet() {
			markup.InlineKeyboard[0][1] = tg.NewInlineKeyboardButton(
				ctx.T().Sprintf("sticker__button_remove-set"), common.DataRemoveSet,
			)
		}
	}

	reply := tg.NewMessage(ctx.Message.Chat.ID, ctx.T().Sprintf("sticker__text"))
	reply.ReplyToMessageID = ctx.Message.ID
	reply.ReplyMarkup = markup
	_, err = ctx.SendMessage(reply)

	return ctx.Error(err)
}
