package middleware

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func Hacktober(bot *tg.Bot, month time.Month) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		if !ctx.IsMessage() {
			return next(ctx)
		}

		lastSeen := time.Unix(ctx.User.LastSeen, 0)
		date := ctx.Message.Time()
		before := time.Date(date.Year(), month, 1, 0, 0, 0, 0, time.UTC)
		// NOTE(toby3d): not November 1, use October 31
		after := before.AddDate(0, 1, 0).Add(-1 * 24 * time.Hour)
		if date.Before(before) || date.After(after) || lastSeen.After(before) {
			return next(ctx)
		}

		// NOTE(toby3d): do this middleware only after sending all previous messages
		if err = next(ctx); err != nil {
			return err
		}

		reply := tg.NewMessage(ctx.Message.Chat.ID, ctx.Printer.Sprintf("🕺 HacktoberFest is here!\n\nIf you are a beginner or already an experienced golang-developer, now is a great time to help improve the quality of the code of this bot. Choose issue to your taste and offer your PR!"))
		reply.DisableNotification = false
		reply.DisableWebPagePreview = false
		reply.ParseMode = tg.StyleMarkdown
		reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButtonURL(
			ctx.Printer.Sprintf("🔧 Let's hack!"), "https://gitlab.com/toby3d/mypackbot/issues?label_name%5B%5D=hacktoberfest",
		)))

		if _, err = bot.SendMessage(reply); err != nil {
			return err
		}

		return next(ctx)
	}
}
