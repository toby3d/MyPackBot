package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

// IsCommand defines actions for commands only
func (h *Handler) IsCommand(ctx *model.Context) (err error) {
	if !ctx.Request.Message.IsCommand() {
		return nil
	}

	switch ctx.Request.Message.Command() {
	case tg.CommandStart:
		err = h.CommandStart(ctx)
	case tg.CommandHelp:
		err = h.CommandHelp(ctx)
	case tg.CommandSettings:
		err = h.CommandSettings(ctx)
	case "ping":
		err = h.sendMessage(ctx, "🏓")
	case "add":
		err = h.CommandAdd(ctx)
	case "del":
		err = h.CommandDel(ctx)
	case "edit":
		err = h.CommandEdit(ctx)
	case "addsticker":
		err = h.CommandAddSticker(ctx)
	case "addpack":
		err = h.CommandAddSet(ctx)
	case "delsticker":
		err = h.CommandDelSticker(ctx)
	case "delpack":
		err = h.CommandDelSet(ctx)
	case "reset":
		err = h.CommandUnknown(ctx)
	case "cancel":
		err = h.CommandUnknown(ctx)
	}

	return err
}

// CommandPing send common ping message.
func (h *Handler) CommandPing(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, "🏓")
}

// CommandStart send common welcome message.
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandStart(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, "start__text "+ctx.Request.Message.From.FullName())
}

// CommandHelp send common message with list of available commands
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandHelp(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, "help__text")
}

// CommandSettings send common message with settings buttons
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandSettings(ctx *model.Context) (err error) {
	return h.CommandUnknown(ctx)
}

// CommandUnknown reply common error message to any unkwnon commands.
func (h *Handler) CommandUnknown(ctx *model.Context) (err error) {
	return h.sendMessage(ctx, "unknown-command__text")
}

func (h *Handler) sendMessage(ctx *model.Context, text string) (err error) {
	reply := tg.NewMessage(ctx.Request.Message.Chat.ID, text)
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Request.Message.ID
	_, err = ctx.SendMessage(reply)

	return err
}

func (h *Handler) CommandAdd(ctx *model.Context) (err error) {
	switch {
	case ctx.Photo != nil:
		err = h.CommandAddPhoto(ctx)
	case ctx.Sticker != nil:
		err = h.CommandAddSticker(ctx)
	}

	h.sendMessage(ctx, "added")

	return err
}

func (h *Handler) CommandEdit(ctx *model.Context) (err error) {
	switch {
	case ctx.Photo != nil:
		err = h.CommandEditPhoto(ctx)
	case ctx.Sticker != nil:
		err = h.CommandEditSticker(ctx)
	}

	h.sendMessage(ctx, "edited")

	return err
}

func (h *Handler) CommandDel(ctx *model.Context) (err error) {
	switch {
	case ctx.Photo != nil:
		err = h.CommandDelPhoto(ctx)
	case ctx.Sticker != nil:
		err = h.CommandDelSticker(ctx)
	}

	h.sendMessage(ctx, "deleted")

	return err
}
