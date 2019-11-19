package middleware

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
)

func UpdateLastSeen(us store.UsersManager) Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) error {
		timeStamp := time.Now().UTC().Unix()

		switch {
		case ctx.IsMessage():
			timeStamp = ctx.Message.Date
		case ctx.IsCallbackQuery():
			timeStamp = ctx.CallbackQuery.Message.Date
		}

		if time.Unix(ctx.User.LastSeen, 0).After(time.Unix(timeStamp, 0).Add(-1 * time.Hour)) {
			return next(ctx)
		}

		ctx.User.LastSeen = timeStamp
		if err := us.Update(ctx.User); err != nil {
			return err
		}

		return next(ctx)
	}
}
