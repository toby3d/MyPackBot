package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
)

type (
	Interceptor   func(*model.Context, model.UpdateFunc) error
	UpdateHandler model.UpdateFunc
	Chain         []Interceptor
)

func (count UpdateHandler) Intercept(middleware Interceptor) UpdateHandler {
	return func(ctx *model.Context) error { return middleware(ctx, model.UpdateFunc(count)) }
}

func (chain Chain) UpdateHandler(handler model.UpdateFunc) model.UpdateFunc {
	current := UpdateHandler(handler)

	for i := len(chain) - 1; i >= 0; i-- {
		m := chain[i]
		current = current.Intercept(m)
	}

	return model.UpdateFunc(current)
}
