package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func AcquirePrinter() Interceptor {
	matcher := language.NewMatcher(message.DefaultCatalog.Languages())

	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		tag, err := language.Parse(ctx.User.LanguageCode)
		if err != nil {
			tag = language.English
		}

		tag, _, _ = matcher.Match(tag)
		p := message.NewPrinter(tag)
		ctx.Set("printer", p)

		return next(ctx)
	}
}
