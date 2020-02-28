package handler

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) GetPhotoKeyboard(ctx *model.Context) *tg.InlineKeyboardMarkup {
	if ctx.Photo == nil {
		return nil
	}

	p := ctx.Get("printer").(*message.Printer)
	userPhoto := h.usersPhotos.Get(&model.UserPhoto{
		UserID:  ctx.User.ID,
		PhotoID: ctx.Photo.ID,
	})
	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("📥 Import this photo"), common.DataAdd),
	))

	if userPhoto != nil {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(p.Sprintf("🔥 Remove this photo"), common.DataDel),
		))
	}

	return &markup
}

func (h *Handler) CommandAddPhoto(ctx *model.Context) (err error) {
	if ctx.Photo == nil || ctx.HasPhoto {
		return nil
	}

	now := time.Now().UTC().Unix()
	userPhoto := new(model.UserPhoto)
	userPhoto.CreatedAt, userPhoto.UpdatedAt = now, now
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

	if !ctx.Request.Message.HasCommandArgument() {
		p := ctx.Get("printer").(*message.Printer)
		reply := tg.NewMessage(int64(ctx.User.ID), p.Sprintf("💡 Add any text and/or emoji(s) as an argument "+
			"of this command to change its search properties."))
		reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
		reply.ParseMode = tg.ParseModeMarkdownV2
		reply.ReplyToMessageID = ctx.Request.Message.ID

		_, err = ctx.SendMessage(reply)

		return err
	}

	if !ctx.HasPhoto {
		if err = h.CommandAddPhoto(ctx); err != nil {
			return err
		}
	}

	return h.usersPhotos.Update(&model.UserPhoto{
		UserID:    ctx.User.ID,
		PhotoID:   ctx.Photo.ID,
		UpdatedAt: ctx.Request.Message.Date,
		Query:     ctx.Request.Message.CommandArgument(),
	})
}

func (h *Handler) CommandDelPhoto(ctx *model.Context) (err error) {
	if ctx.Photo == nil || !ctx.HasPhoto {
		return nil
	}

	return h.usersPhotos.Remove(&model.UserPhoto{
		UserID:  ctx.User.ID,
		PhotoID: ctx.Photo.ID,
	})
}
