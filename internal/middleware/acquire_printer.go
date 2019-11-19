package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func AcquirePrinter() Interceptor {
	matcher := message.DefaultCatalog.Matcher()

	return func(ctx *model.Context, next model.UpdateFunc) error {
		tag, _, _ := matcher.Match(language.Make(ctx.User.LanguageCode))
		ctx.Printer = message.NewPrinter(tag)
		return next(ctx)
	}
}
