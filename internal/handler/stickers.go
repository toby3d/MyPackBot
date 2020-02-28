package handler

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) GetStickerKeyboard(ctx *model.Context) *tg.InlineKeyboardMarkup {
	if ctx.Sticker == nil {
		return nil
	}

	p := ctx.Get("printer").(*message.Printer)
	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton("🔥 Remove this sticker", common.DataDel),
	))

	if (ctx.Request.IsCallbackQuery() && ctx.Request.CallbackQuery.Data == common.DataDelSet) || !ctx.HasSticker {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(p.Sprintf("📥 Import this sticker"), common.DataAdd),
		))
	}

	if ctx.Sticker.InSet() {
		setName, _ := ctx.Get("set_name").(string)
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(p.Sprintf("📥 Import %s set", setName), common.DataAddSet),
		))

		if ctx.Request.IsCallbackQuery() && ctx.Request.CallbackQuery.Data == common.DataAddSet || ctx.HasSet {
			markup.InlineKeyboard[1][0] = tg.NewInlineKeyboardButton(
				p.Sprintf("🔥 Remove %s set", setName), common.DataDelSet,
			)
		}
	}

	return &markup
}

// CommandAddSticker import single Sticker by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandAddSticker(ctx *model.Context) (err error) {
	if ctx.HasSticker || ctx.Sticker == nil {
		return nil
	}

	now := time.Now().UTC().Unix()
	userSticker := new(model.UserSticker)
	userSticker.CreatedAt, userSticker.UpdatedAt = now, now
	userSticker.Query = ctx.Sticker.Emoji
	userSticker.StickerID = ctx.Sticker.ID
	userSticker.UserID = ctx.User.ID

	if ctx.Request.IsMessage() && ctx.Request.Message.HasCommandArgument() {
		userSticker.Query = ctx.Request.Message.CommandArgument()
	}

	return h.usersStickers.Add(userSticker)
}

// CommandAddPack import whole Sticker pack by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandAddSet(ctx *model.Context) (err error) {
	if ctx.Sticker == nil || !ctx.Sticker.InSet() {
		return nil
	}

	return h.usersStickers.AddSet(ctx.User.ID, ctx.Sticker.SetName)
}

func (h *Handler) CommandEditSticker(ctx *model.Context) (err error) {
	if ctx.Sticker == nil {
		return nil
	}

	if !ctx.Request.Message.HasCommandArgument() {
		p := ctx.Get("printer").(*message.Printer)
		reply := tg.NewMessage(int64(ctx.User.ID), p.Sprintf("💡 Add any text and/or emoji(s) as an "+
			"argument of this command to change its search properties."))
		reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
		reply.ParseMode = tg.ParseModeMarkdownV2
		reply.ReplyToMessageID = ctx.Request.Message.ID

		_, err = ctx.SendMessage(reply)

		return err
	}

	if !ctx.HasSticker {
		if err = h.CommandAddSticker(ctx); err != nil {
			return err
		}
	}

	return h.usersStickers.Update(&model.UserSticker{
		UserID:    ctx.User.ID,
		StickerID: ctx.Sticker.ID,
		UpdatedAt: ctx.Request.Message.Date,
		Query:     ctx.Request.Message.CommandArgument(),
	})
}

// CommandDelSticker remove single Sticker by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandDelSticker(ctx *model.Context) (err error) {
	if ctx.Sticker == nil || !ctx.HasSticker {
		return nil
	}

	return h.usersStickers.Remove(&model.UserSticker{
		UserID:    ctx.User.ID,
		StickerID: ctx.Sticker.ID,
	})
}

// CommandDelPack remove whole Sticker pack by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandDelSet(ctx *model.Context) (err error) {
	if ctx.Sticker == nil || !ctx.Sticker.InSet() {
		return nil
	}

	return h.usersStickers.RemoveSet(ctx.User.ID, ctx.Sticker.SetName)
}
