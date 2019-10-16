package middleware

import (
	"context"

	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

type (
	Interceptor   func(context.Context, *tg.Update, model.UpdateFunc) error
	UpdateHandler model.UpdateFunc
	Chain         []Interceptor
)

func (count UpdateHandler) Intercept(middleware Interceptor) UpdateHandler {
	return func(ctx context.Context, upd *tg.Update) error {
		return middleware(ctx, upd, model.UpdateFunc(count))
	}
}

func (chain Chain) UpdateHandler(handler model.UpdateFunc) model.UpdateFunc {
	current := UpdateHandler(handler)
	for i := len(chain) - 1; i >= 0; i-- {
		m := chain[i]
		current = current.Intercept(m)
	}
	return model.UpdateFunc(current)
}
