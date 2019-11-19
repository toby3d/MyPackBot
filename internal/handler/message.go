package handler

import (
	"strings"

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

	return err
}

func (h *Handler) IsCommand(ctx *model.Context) (err error) {
	switch {
	case ctx.Message.IsCommandEqual(tg.CommandStart):
		err = h.CommandStart(ctx)
	case ctx.Message.IsCommandEqual(tg.CommandHelp):
		err = h.CommandHelp(ctx)
	case ctx.Message.IsCommandEqual(tg.CommandSettings),
		ctx.Message.IsCommandEqual("addSticker"),
		ctx.Message.IsCommandEqual("addPack"),
		ctx.Message.IsCommandEqual("delSticker"),
		ctx.Message.IsCommandEqual("delPack"),
		ctx.Message.IsCommandEqual("reset"),
		ctx.Message.IsCommandEqual("cancel"):
		fallthrough
	default:
		err = h.CommandUnknown(ctx)
	}

	return err
}

func (h *Handler) CommandStart(ctx *model.Context) (err error) {
	reply := tg.NewMessage(ctx.Message.Chat.ID, ctx.T().Sprintf("start__text", ctx.Message.From.FullName()))
	reply.ReplyToMessageID = ctx.Message.ID
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)

	_, err = ctx.SendMessage(reply)

	return err
}

func (h *Handler) CommandHelp(ctx *model.Context) (err error) {
	reply := tg.NewMessage(ctx.Message.Chat.ID, ctx.T().Sprintf("help__text"))
	reply.ReplyToMessageID = ctx.Message.ID
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)

	_, err = ctx.SendMessage(reply)

	return err
}

func (h *Handler) CommandUnknown(ctx *model.Context) (err error) {
	reply := tg.NewMessage(ctx.Message.Chat.ID, ctx.T().Sprintf("unknown-command__text"))
	reply.ReplyToMessageID = ctx.Message.ID
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)

	_, err = ctx.SendMessage(reply)

	return err
}

func (h *Handler) IsSticker(ctx *model.Context) error {
	us, err := h.store.GetSticker(ctx.User, ctx.Sticker)
	if err != nil {
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(ctx.T().Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))

	if !strings.EqualFold(ctx.Sticker.SetName, common.SetNameUploaded) {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(ctx.T().Sprintf("sticker__button_add-set"), common.DataAddSet),
		)
	}

	if us != nil {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButton(
			ctx.T().Sprintf("sticker__button_remove-single"), common.DataRemoveSticker,
		)))

		if ctx.Sticker.SetName != "" {
			markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], tg.NewInlineKeyboardButton(
				ctx.T().Sprintf("sticker__button_remove-set"), common.DataRemoveSet,
			))
		}
	}

	reply := tg.NewMessage(ctx.Message.Chat.ID, ctx.T().Sprintf("sticker__text"))
	reply.ReplyToMessageID = ctx.Message.ID
	reply.ReplyMarkup = markup

	_, err = ctx.SendMessage(reply)

	return err
}
