package handler

import (
	"fmt"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) GetStickerKeyboard(ctx *model.Context) *tg.InlineKeyboardMarkup {
	if ctx.Sticker == nil {
		return nil
	}

	ctx.UserSticker = h.usersStickers.Get(&model.UserSticker{
		UserID:    ctx.User.ID,
		StickerID: ctx.Sticker.ID,
	})

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton("del", common.DataDel),
	))

	if (ctx.Request.IsCallbackQuery() && ctx.Request.CallbackQuery.Data == common.DataDelSet) ||
		ctx.UserSticker == nil {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton("add", common.DataAdd),
		))
	}

	if ctx.Sticker.InSet() {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton("add set", common.DataAddSet),
		))

		if (ctx.Request.IsCallbackQuery() && ctx.Request.CallbackQuery.Data == common.DataAddSet) ||
			ctx.UserSticker != nil {
			markup.InlineKeyboard[1][0] = tg.NewInlineKeyboardButton("del set", common.DataDelSet)
		}
	}

	return markup
}

// CommandAddSticker import single Sticker by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandAddSticker(ctx *model.Context) (err error) {
	if ctx.Sticker == nil || ctx.UserSticker != nil {
		return nil
	}

	userSticker := new(model.UserSticker)
	userSticker.UserID = ctx.User.ID
	userSticker.StickerID = ctx.Sticker.ID

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

	if ctx.UserSticker == nil {
		return h.CommandAddSticker(ctx)
	}

	if !ctx.Request.Message.HasCommandArgument() {
		query := ctx.UserSticker.Query
		if query == "" {
			query = h.stickers.Get(ctx.UserSticker.StickerID).Emoji
		}

		return h.sendMessage(ctx, fmt.Sprintln("current query:", query))
	}

	ctx.UserSticker.UpdatedAt = ctx.Request.Message.Date
	ctx.UserSticker.Query = ctx.Request.Message.CommandArgument()

	return h.usersStickers.Update(ctx.UserSticker)
}

// CommandDelSticker remove single Sticker by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandDelSticker(ctx *model.Context) (err error) {
	if ctx.Sticker == nil || ctx.UserSticker == nil {
		return nil
	}

	return h.usersStickers.Remove(ctx.UserSticker)
}

// CommandDelPack remove whole Sticker pack by ReplyMessage.
// NOTE(toby3d): DEPRECATED, used for backward compatibility
func (h *Handler) CommandDelSet(ctx *model.Context) (err error) {
	if ctx.Sticker == nil || !ctx.Sticker.InSet() {
		return nil
	}

	return h.usersStickers.RemoveSet(ctx.User.ID, ctx.Sticker.SetName)
}
