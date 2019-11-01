package handler

import (
	"context"
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) IsMessage(ctx context.Context, msg *tg.Message) (err error) {
	switch {
	case msg.IsCommand():
		err = h.IsCommand(ctx, msg)
	case msg.IsSticker():
		err = h.IsSticker(ctx, msg)
	}

	return err
}

func (h *Handler) IsCommand(ctx context.Context, msg *tg.Message) (err error) {
	switch {
	case msg.IsCommandEqual(tg.CommandStart):
		err = h.CommandStart(ctx, msg)
	case msg.IsCommandEqual(tg.CommandHelp):
		err = h.CommandHelp(ctx, msg)
	case msg.IsCommandEqual(tg.CommandSettings):
		fallthrough
	default:
		err = h.CommandUnknown(ctx, msg)
	}

	return err
}

func (h *Handler) CommandStart(ctx context.Context, msg *tg.Message) (err error) {
	p, _ := ctx.Value(common.ContextPrinter).(*message.Printer)

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("start__text", msg.From.FullName()))
	reply.ReplyToMessageID = msg.ID
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)

	_, err = h.bot.SendMessage(reply)

	return err
}

func (h *Handler) CommandHelp(ctx context.Context, msg *tg.Message) (err error) {
	p, _ := ctx.Value(common.ContextPrinter).(*message.Printer)

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("help__text"))
	reply.ReplyToMessageID = msg.ID

	_, err = h.bot.SendMessage(reply)

	return err
}

func (h *Handler) CommandUnknown(ctx context.Context, msg *tg.Message) (err error) {
	p, _ := ctx.Value(common.ContextPrinter).(*message.Printer)

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("unknown-command__text"))
	reply.ReplyToMessageID = msg.ID

	_, err = h.bot.SendMessage(reply)

	return err
}

func (h *Handler) IsSticker(ctx context.Context, msg *tg.Message) error {
	u, _ := ctx.Value(common.ContextUser).(*model.User)
	p, _ := ctx.Value(common.ContextPrinter).(*message.Printer)
	s, _ := ctx.Value(common.ContextSticker).(*model.Sticker)

	us, err := h.store.GetSticker(u, s)
	if err != nil {
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))

	if !strings.EqualFold(s.SetName, common.SetNameUploaded) {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-set"), common.DataAddSet),
		)
	}

	if us != nil {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButton(
			p.Sprintf("sticker__button_remove-single"),
			common.DataRemoveSticker,
		)))

		if s.SetName != "" {
			markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], tg.NewInlineKeyboardButton(
				p.Sprintf("sticker__button_remove-set"),
				common.DataRemoveSet,
			))
		}
	}

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("sticker__text"))
	reply.ReplyToMessageID = msg.ID
	reply.ReplyMarkup = markup

	_, err = h.bot.SendMessage(reply)

	return err
}
