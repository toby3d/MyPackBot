package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func AcquirePrinter() Interceptor {
	matcher := language.NewMatcher(message.DefaultCatalog.Languages())
	return func(ctx *model.Context, next model.UpdateFunc) (err error) {
		tag, _, _ := matcher.Match(language.MustParse(ctx.User.LanguageCode))
		p := message.NewPrinter(tag)
		ctx.Set("printer", p)
		return next(ctx)
	}
}
