package middleware

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
)

func UpdateLastSeen(us users.Manager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		timeStamp := time.Now().UTC().Unix()

		switch {
		case ctx.Request.IsMessage():
			timeStamp = ctx.Request.Message.Date
		case ctx.Request.IsCallbackQuery():
			timeStamp = ctx.Request.CallbackQuery.Message.Date
		}

		if time.Unix(ctx.User.LastSeen, 0).After(time.Unix(timeStamp, 0).Add(-1 * time.Hour)) {
			return next(ctx)
		}

		ctx.User.LastSeen = timeStamp
		if err = us.Update(ctx.User); err != nil {
			return err
		}

		return next(ctx)
	}
}
