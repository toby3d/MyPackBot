package handler

import (
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
	ctx.UserPhoto = h.usersPhotos.Get(&model.UserPhoto{
		UserID:  ctx.User.ID,
		PhotoID: ctx.Photo.ID,
	})

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("📥 Import this photo"), common.DataAdd),
	))

	if ctx.UserPhoto != nil {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(p.Sprintf("🔥 Remove this photo"), common.DataDel),
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

	if !ctx.Request.Message.HasCommandArgument() {
		p := ctx.Get("printer").(*message.Printer)
		reply := tg.NewMessage(ctx.User.UserID, p.Sprintf("💡 Add any text and/or emoji(s) as an argument of "+
			"this command to change its search properties."))
		reply.ReplyMarkup = tg.NewReplyKeyboardRemove(false)
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyToMessageID = ctx.Request.Message.ID

		_, err = ctx.SendMessage(reply)

		return err
	}

	if ctx.UserPhoto == nil {
		if err = h.CommandAddPhoto(ctx); err != nil {
			return err
		}
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
