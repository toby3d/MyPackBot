package middleware

import (
	"context"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func Birthday(bot *tg.Bot, us store.UsersManager, bday time.Time) Interceptor {
	return func(ctx context.Context, update *tg.Update, next model.UpdateFunc) (err error) {
		if !update.IsMessage() {
			return next(ctx, update)
		}

		u, _ := ctx.Value("user").(*model.User)
		lastSeen := time.Unix(u.LastSeen, 0)
		date := update.Message.Time()
		before := time.Date(date.Year(), bday.Month(), bday.Day(), 0, 0, 0, 0, time.UTC)
		after := before.AddDate(0, 0, 7)
		if date.Before(before) || date.After(after) || lastSeen.After(before) {
			return next(ctx, update)
		}
		// NOTE(toby3d): do this middleware only after sending all previous messages
		if err = next(ctx, update); err != nil {
			return err
		}

		p, _ := ctx.Value("printer").(*message.Printer)
		reply := tg.NewMessage(update.Message.Chat.ID, p.Sprintf("birthday__message_text"))
		reply.DisableNotification = false
		reply.DisableWebPagePreview = false
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonURL(
			p.Sprintf("birthday__button_text-donate"), "https://toby3d.me/donate",
		)))

		if _, err = bot.SendMessage(reply); err != nil {
			return err
		}

		u.LastSeen = date.Unix()
		if err = us.Update(u); err != nil {
			return err
		}

		return next(ctx, update)
	}
}
