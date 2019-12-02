package middleware

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
)

func AcquireUser(us users.Manager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		ctx.User = new(model.User)

		switch {
		case ctx.Request.IsMessage():
			ctx.User.UserID = int64(ctx.Request.Message.From.ID)
			ctx.User.CreatedAt = ctx.Request.Message.Date
			ctx.User.UpdatedAt = ctx.Request.Message.Date
			ctx.User.LanguageCode = ctx.Request.Message.From.LanguageCode
			ctx.User.LastSeen = ctx.Request.Message.Date
		case ctx.Request.IsInlineQuery():
			now := time.Now().UTC().Unix()
			ctx.User.UserID = int64(ctx.Request.InlineQuery.From.ID)
			ctx.User.CreatedAt = now
			ctx.User.UpdatedAt = now
			ctx.User.LanguageCode = ctx.Request.InlineQuery.From.LanguageCode
			ctx.User.LastSeen = now
		case ctx.Request.IsCallbackQuery():
			ctx.User.UserID = int64(ctx.Request.CallbackQuery.From.ID)
			ctx.User.CreatedAt = ctx.Request.CallbackQuery.Message.Date
			ctx.User.UpdatedAt = ctx.Request.CallbackQuery.Message.Date
			ctx.User.LanguageCode = ctx.Request.CallbackQuery.From.LanguageCode
			ctx.User.LastSeen = ctx.Request.CallbackQuery.Message.Date
		}

		if ctx.User, err = us.GetOrCreate(ctx.User); err != nil {
			return err
		}

		return next(ctx)
	}
}
