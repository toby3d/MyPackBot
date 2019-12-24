package handler

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
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
		err = h.CommandUnknown(ctx)
	case "ping":
		err = h.CommandPing(ctx)
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
	case "reset", "cancel":
		err = h.CommandUnknown(ctx)
	}

	return err
}

// CommandPing send common ping message.
func (h *Handler) CommandPing(ctx *model.Context) (err error) {
	reply := tg.NewMessage(ctx.User.UserID, "🏓")
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Request.Message.ID
	_, err = ctx.SendMessage(reply)
	return err
}

// CommandStart send common welcome message.
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandStart(ctx *model.Context) (err error) {
	p := ctx.Get("printer").(*message.Printer)
	reply := tg.NewMessage(ctx.User.UserID, p.Sprintf("👋 Hi %s, I'm %s!\nThanks to me, you can collect almost any"+
		" media content in Telegram without any limits, in any chat via inline mode.",
		ctx.Request.Message.From.FullName(), ctx.FullName()))
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Request.Message.ID
	_, err = ctx.SendMessage(reply)
	return err
}

// CommandHelp send common message with list of available commands
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandHelp(ctx *model.Context) (err error) {
	p := ctx.Get("printer").(*message.Printer)
	reply := tg.NewMessage(ctx.User.UserID, p.Sprintf("🤖 Here is a list of commands that I understand, some of"+
		" them [may] or (should) contain an argument:\n/start - start all over again\n/help [other command] "+
		"- get a list of available commands or help and a demonstration of a specific command\n/add [query] "+
		"- add media from reply to your collection [with custom search query]\n/edit (query) - change query "+
		"to reply media\n/del - remove reply media from your collection"))
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Request.Message.ID
	_, err = ctx.SendMessage(reply)
	return err
}

// CommandSettings send common message with settings buttons
// NOTE(toby3d): REQUIRED by Telegram Bot API platform
func (h *Handler) CommandSettings(ctx *model.Context) (err error) {
	return nil
}

// CommandUnknown reply common error message to any unkwnon commands.
func (h *Handler) CommandUnknown(ctx *model.Context) (err error) {
	return nil
}

func (h *Handler) CommandAdd(ctx *model.Context) (err error) {
	switch {
	case ctx.Photo != nil:
		err = h.CommandAddPhoto(ctx)
	case ctx.Sticker != nil:
		err = h.CommandAddSticker(ctx)
	default:
		return nil
	}
	if err != nil {
		return err
	}

	if !ctx.Request.Message.HasCommandArgument() {
		return err
	}

	p := ctx.Get("printer").(*message.Printer)
	reply := tg.NewMessage(ctx.User.UserID, p.Sprintf("👍 Imported!"))
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Request.Message.ID
	_, err = ctx.SendMessage(reply)

	return err
}

func (h *Handler) CommandEdit(ctx *model.Context) (err error) {
	switch {
	case ctx.Photo != nil:
		err = h.CommandEditPhoto(ctx)
	case ctx.Sticker != nil:
		err = h.CommandEditSticker(ctx)
	default:
		return nil
	}
	if err != nil {
		return err
	}

	if !ctx.Request.Message.HasCommandArgument() {
		return err
	}

	p := ctx.Get("printer").(*message.Printer)
	reply := tg.NewMessage(ctx.User.UserID, p.Sprintf("👍 Updated!"))
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Request.Message.ID
	_, err = ctx.SendMessage(reply)

	return err
}

func (h *Handler) CommandDel(ctx *model.Context) (err error) {
	switch {
	case ctx.Photo != nil:
		err = h.CommandDelPhoto(ctx)
	case ctx.Sticker != nil:
		err = h.CommandDelSticker(ctx)
	default:
		return nil
	}
	if err != nil {
		return err
	}

	if !ctx.Request.Message.HasCommandArgument() {
		return err
	}

	p := ctx.Get("printer").(*message.Printer)
	reply := tg.NewMessage(ctx.User.UserID, p.Sprintf("👍 Removed!"))
	reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
	reply.ReplyToMessageID = ctx.Request.Message.ID
	_, err = ctx.SendMessage(reply)

	return err
}
