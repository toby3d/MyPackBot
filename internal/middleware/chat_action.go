package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func ChatAction() Interceptor {
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		if !ctx.Request.IsMessage() || !ctx.Request.Message.Chat.IsPrivate() {
			return next(ctx)
		}

		go func() { _, _ = ctx.SendChatAction(ctx.Request.Message.Chat.ID, tg.ActionTyping) }()
		return next(ctx)
	}
}
