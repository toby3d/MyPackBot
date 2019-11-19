package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func ChatAction() Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		if !ctx.IsMessage() || !ctx.Message.Chat.IsPrivate() {
			return next(ctx)
		}

		if _, err = ctx.SendChatAction(ctx.Message.Chat.ID, tg.ActionTyping); err != nil {
			return ctx.Error(err)
		}

		return next(ctx)
	}
}
