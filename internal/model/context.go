package model

import (
	"context"

	tg "gitlab.com/toby3d/telegram"
)

type (
	UpdateFunc func(*Context) error

	Context struct {
		*tg.Bot
		Request *tg.Update

		User       *User
		Sticker    *Sticker
		HasSticker bool
		HasSet     bool
		Photo      *Photo
		HasPhoto   bool

		userValues context.Context
	}

	contextKey string
)

func (ctx *Context) Set(key string, val interface{}) {
	if ctx.userValues == nil {
		ctx.userValues = context.Background()
	}

	ctx.userValues = context.WithValue(ctx.userValues, contextKey(key), val)
}

func (ctx *Context) Get(key string) interface{} {
	if ctx.userValues == nil {
		ctx.userValues = context.Background()
	}

	return ctx.userValues.Value(contextKey(key))
}
