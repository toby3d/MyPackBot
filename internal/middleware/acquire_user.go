package middleware

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
)

func AcquireUser(us users.Manager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		switch {
		case ctx.IsMessage():
			ctx.User.ID = ctx.Message.From.ID
			ctx.User.CreatedAt = ctx.Message.Date
			ctx.User.UpdatedAt = ctx.Message.Date
			ctx.User.LanguageCode = ctx.Message.From.LanguageCode
			ctx.User.LastSeen = ctx.Message.Date
		case ctx.IsInlineQuery():
			now := time.Now().UTC().Unix()
			ctx.User.ID = ctx.InlineQuery.From.ID
			ctx.User.CreatedAt = now
			ctx.User.UpdatedAt = now
			ctx.User.LanguageCode = ctx.InlineQuery.From.LanguageCode
			ctx.User.LastSeen = now
		case ctx.IsCallbackQuery():
			ctx.User.ID = ctx.CallbackQuery.From.ID
			ctx.User.CreatedAt = ctx.CallbackQuery.Message.Date
			ctx.User.UpdatedAt = ctx.CallbackQuery.Message.Date
			ctx.User.LanguageCode = ctx.CallbackQuery.From.LanguageCode
			ctx.User.LastSeen = ctx.CallbackQuery.Message.Date
		}

		if ctx.User, err = us.GetOrCreate(ctx.User); err != nil {
			return err
		}

		return next(ctx)
	}
}
