package middleware

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func Birthday(bday time.Time) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		if !ctx.Request.IsMessage() {
			return next(ctx)
		}

		lastSeen := time.Unix(ctx.User.LastSeen, 0)
		date := ctx.Request.Message.Time()
		before := time.Date(date.Year(), bday.Month(), bday.Day(), 0, 0, 0, 0, time.UTC)
		after := before.AddDate(0, 0, 7)
		if date.Before(before) || date.After(after) || lastSeen.After(before) {
			return next(ctx)
		}

		// NOTE(toby3d): do this middleware only after sending all previous messages
		if err = next(ctx); err != nil {
			return err
		}

		p := ctx.Get("printer").(*message.Printer)
		reply := tg.NewMessage(ctx.Request.Message.Chat.ID, p.Sprintf("🥳 4 November? It's a @toby3d birthday!"+
			"\nIf you like this bot, then why not send him a congratulation along with a small gift? This"+
			" will make him incredibly happy!"))
		if date.After(bday.AddDate(0, 0, 1)) {
			reply.Text = p.Sprintf("☺️ Oh, you missed @toby3d birthday on November 4th!\nIf you like this" +
				" bot, why not send him some birthday greetings and a little birthday gift? It is " +
				"not yet too late to make him happy!")
		}
		reply.DisableNotification = false
		reply.DisableWebPagePreview = false
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonURL(
				p.Sprintf("💸 Make a donation"), "https://toby3d.me/donate",
			)), tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonURL(
				p.Sprintf("🤝 Use referral links"), "https://toby3d.me/referrals",
			)),
		)

		if _, err = ctx.SendMessage(reply); err != nil {
			return err
		}

		return next(ctx)
	}
}
