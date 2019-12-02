package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) IsMessage(ctx *model.Context) (err error) {
	switch {
	case ctx.Request.Message.IsCommand():
		err = h.IsCommand(ctx)
	case ctx.Request.Message.IsSticker():
		_, err = ctx.SendMessage(&tg.SendMessageParameters{
			ChatID:           ctx.Request.Message.Chat.ID,
			ReplyMarkup:      h.GetStickerKeyboard(ctx),
			ReplyToMessageID: ctx.Request.Message.ID,
			Text:             tg.TypeSticker,
		})
	case ctx.Request.Message.IsPhoto():
		_, err = ctx.SendMessage(&tg.SendMessageParameters{
			ChatID:           ctx.Request.Message.Chat.ID,
			ReplyMarkup:      h.GetPhotoKeyboard(ctx),
			ReplyToMessageID: ctx.Request.Message.ID,
			Text:             tg.TypePhoto,
		})
	}

	return err
}
