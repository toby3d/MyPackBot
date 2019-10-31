package handler

import (
	"context"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) isMessage(ctx context.Context, msg *tg.Message) (err error) {
	switch {
	case msg.IsCommand():
		err = h.isCommand(ctx, msg)
	case msg.IsSticker():
		err = h.isSticker(ctx, msg)
	}
	return err
}

func (h *Handler) isCommand(ctx context.Context, msg *tg.Message) (err error) {
	switch {
	case msg.IsCommandEqual(tg.CommandStart):
		err = h.commandStart(ctx, msg)
	case msg.IsCommandEqual(tg.CommandHelp):
		err = h.commandHelp(ctx, msg)
	case msg.IsCommandEqual(tg.CommandSettings):
		fallthrough
	default:
		err = h.commandUnknown(ctx, msg)
	}
	return err
}

func (h *Handler) commandStart(ctx context.Context, msg *tg.Message) (err error) {
	p, _ := ctx.Value("printer").(*message.Printer)
	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("start__text", msg.From.FullName()))
	reply.ReplyToMessageID = msg.ID
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	_, err = h.bot.SendMessage(reply)
	return err
}

func (h *Handler) commandHelp(ctx context.Context, msg *tg.Message) (err error) {
	p, _ := ctx.Value("printer").(*message.Printer)
	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("help__text"))
	reply.ReplyToMessageID = msg.ID
	_, err = h.bot.SendMessage(reply)
	return err
}

func (h *Handler) commandUnknown(ctx context.Context, msg *tg.Message) (err error) {
	p, _ := ctx.Value("printer").(*message.Printer)
	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("unknown-command__text"))
	reply.ReplyToMessageID = msg.ID
	_, err = h.bot.SendMessage(reply)
	return err
}

func (h *Handler) isSticker(ctx context.Context, msg *tg.Message) error {
	u, _ := ctx.Value("user").(*model.User)
	p, _ := ctx.Value("printer").(*message.Printer)
	s, _ := ctx.Value("sticker").(*model.Sticker)

	us, err := h.store.GetSticker(u, s)
	if err != nil {
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-single"), common.DataAddSticker),
	))
	if s.SetName != "" {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-set"), common.DataAddSet),
		)
	}
	if us.StickerID != "" && us.UserID != 0 {
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
