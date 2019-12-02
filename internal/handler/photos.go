package handler

import (
	"fmt"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) GetPhotoKeyboard(ctx *model.Context) *tg.InlineKeyboardMarkup {
	if ctx.Photo == nil {
		return nil
	}

	ctx.UserPhoto = h.usersPhotos.Get(&model.UserPhoto{
		UserID:  ctx.User.ID,
		PhotoID: ctx.Photo.ID,
	})

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton("add", common.DataAdd),
	))

	if ctx.UserPhoto != nil {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton("del", common.DataDel),
		))
	}

	return markup
}

func (h *Handler) CommandAddPhoto(ctx *model.Context) (err error) {
	if ctx.Photo == nil || ctx.UserPhoto != nil {
		return nil
	}

	userPhoto := new(model.UserPhoto)
	userPhoto.UserID = ctx.User.ID
	userPhoto.PhotoID = ctx.Photo.ID

	if ctx.Request.IsMessage() && ctx.Request.Message.HasCommandArgument() {
		userPhoto.Query = ctx.Request.Message.CommandArgument()
	}

	return h.usersPhotos.Add(userPhoto)
}

func (h *Handler) CommandEditPhoto(ctx *model.Context) (err error) {
	if ctx.Photo == nil {
		return nil
	}

	if ctx.UserPhoto == nil {
		return h.CommandAddPhoto(ctx)
	}

	if !ctx.Request.Message.HasCommandArgument() {
		return h.sendMessage(ctx, fmt.Sprintln("current query:", ctx.UserPhoto.Query))
	}

	ctx.UserPhoto.UpdatedAt = ctx.Request.Message.Date
	ctx.UserPhoto.Query = ctx.Request.Message.CommandArgument()

	return h.usersPhotos.Update(ctx.UserPhoto)
}

func (h *Handler) CommandDelPhoto(ctx *model.Context) (err error) {
	if ctx.Photo == nil || ctx.UserPhoto == nil {
		return nil
	}

	return h.usersPhotos.Remove(ctx.UserPhoto)
}
