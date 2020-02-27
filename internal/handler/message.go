package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) IsMessage(ctx *model.Context) (err error) {
	p := ctx.Get("printer").(*message.Printer)

	reply := tg.NewMessage(ctx.Request.Message.Chat.ID, p.Sprintf("🤔 What to do with this?"))
	reply.ReplyToMessageID = ctx.Request.Message.ID

	switch {
	case ctx.Request.Message.IsCommand():
		err = h.IsCommand(ctx)
	case ctx.Request.Message.IsSticker():
		reply.ReplyMarkup = h.GetStickerKeyboard(ctx)
		_, err = ctx.SendMessage(reply)
	case ctx.Request.Message.IsPhoto():
		reply.ReplyMarkup = h.GetPhotoKeyboard(ctx)
		_, err = ctx.SendMessage(reply)
	}

	return err
}
