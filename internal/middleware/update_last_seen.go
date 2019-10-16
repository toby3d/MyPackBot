package middleware

import (
	"context"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
)

func UpdateLastSeen(us store.UsersManager) Interceptor {
	return func(ctx context.Context, update *tg.Update, next model.UpdateFunc) error {
		timeStamp := time.Now().UTC().Unix()
		if update.IsMessage() {
			timeStamp = update.Message.Date
		}
		u, _ := ctx.Value("user").(*model.User)
		if time.Unix(u.LastSeen, 0).After(time.Unix(timeStamp, 0).Add(-1 * time.Hour)) {
			return next(ctx, update)
		}

		u.LastSeen = timeStamp
		if err := us.Update(u); err != nil {
			return err
		}

		return next(context.WithValue(ctx, "user", u), update)
	}
}
