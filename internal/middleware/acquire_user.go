package middleware

import (
	"context"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
)

func AcquireUser(us store.UsersManager) Interceptor {
	return func(ctx context.Context, update *tg.Update, next model.UpdateFunc) error {
		timeStamp := time.Now().UTC().Unix()

		from := new(tg.User)
		switch {
		case update.IsMessage():
			*from = *update.Message.From
			timeStamp = update.Message.Date
		case update.IsInlineQuery():
			*from = *update.InlineQuery.From
		case update.IsCallbackQuery():
			*from = *update.CallbackQuery.From
		}

		u, err := us.GetOrCreate(&model.User{
			ID:        from.ID,
			CreatedAt: timeStamp,
			UpdatedAt: timeStamp,

			LanguageCode: from.LanguageCode,
			LastSeen:     timeStamp,
		})
		if err != nil {
			return err
		}

		return next(context.WithValue(ctx, "user", u), update)
	}
}
