package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

// IsCommand defines actions for commands only
func (h *Handler) IsCommand(ctx *model.Context) (err error) {
	switch {
	case ctx.Message.IsCommandEqual(tg.CommandStart):
		err = h.CommandStart(ctx)
	case ctx.Message.IsCommandEqual(tg.CommandHelp):
		err = h.CommandHelp(ctx)
	case ctx.Message.IsCommandEqual(tg.CommandSettings):
		err = h.CommandSettings(ctx)
	case ctx.Message.IsCommandEqual(common.CommandPing):
		err = h.CommandPing(ctx)
	case ctx.Message.IsCommandEqual(common.CommandDelSticker):
		err = h.CommandDelSticker(ctx)
	case ctx.Message.IsCommandEqual(common.CommandDelPack):
		err = h.CommandDelPack(ctx)
	case ctx.Message.IsCommandEqual(common.DataAddSticker):
		err = h.CommandAddSticker(ctx)
	case ctx.Message.IsCommandEqual(common.CommandAddPack):
		err = h.CommandAddPack(ctx)
	case ctx.Message.IsCommandEqual(common.CommandReset):
		fallthrough
	case ctx.Message.IsCommandEqual(common.CommandCancel):
		fallthrough
	default:
		err = h.CommandUnknown(ctx)
	}

	return ctx.Error(err)
}

// CommandPing send common ping message.
func (h *Handler) CommandPing(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, "🏓")
}

// CommandStart send common welcome message.
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandStart(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, ctx.T().Sprintf("start__text", ctx.Message.From.FullName()))
}

// CommandHelp send common message with list of available commands
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandHelp(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, ctx.T().Sprintf("help__text"))
}

// CommandSettings send common message with settings buttons
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandSettings(ctx *model.Context) (err error) {
	return h.CommandUnknown(ctx)
}

// CommandAddSticker import single Sticker by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandAddSticker(ctx *model.Context) (err error) {
	if !ctx.Message.IsReply() || !ctx.Message.ReplyToMessage.IsSticker() {
		return nil
	}

	go h.store.AddSticker(ctx.User, ctx.Sticker)

	return h.sendMessage(ctx, ctx.T().Sprintf("addsticker-command__text"))
}

// CommandAddPack import whole Sticker pack by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandAddPack(ctx *model.Context) (err error) {
	if !ctx.Message.IsReply() || !ctx.Message.ReplyToMessage.IsSticker() {
		return nil
	}

	go h.store.AddStickersSet(ctx.User, ctx.Sticker.SetName)

	return h.sendMessage(ctx, ctx.T().Sprintf("addpack-command__text"))
}

// CommandDelSticker remove single Sticker by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandDelSticker(ctx *model.Context) (err error) {
	if !ctx.Message.IsReply() || !ctx.Message.ReplyToMessage.IsSticker() {
		return nil
	}

	go h.store.RemoveSticker(ctx.User, ctx.Sticker)

	return h.sendMessage(ctx, ctx.T().Sprintf("delsticker-command__text"))
}

// CommandDelPack remove whole Sticker pack by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandDelPack(ctx *model.Context) (err error) {
	if !ctx.Message.IsReply() || !ctx.Message.ReplyToMessage.IsSticker() {
		return nil
	}

	go h.store.RemoveStickersSet(ctx.User, ctx.Sticker.SetName)

	return h.sendMessage(ctx, ctx.T().Sprintf("delpack-command__text"))
}

// CommandUnknown reply common error message to any unkwnon commands.
func (h *Handler) CommandUnknown(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, ctx.T().Sprintf("unknown-command__text"))
}

func (h *Handler) sendMessage(ctx *model.Context, text string) (err error) {
	reply := tg.NewMessage(ctx.Message.Chat.ID, text)
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Message.ID
	_, err = ctx.SendMessage(reply)

	return ctx.Error(err)
}
