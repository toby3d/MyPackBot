package middleware

import (
	"context"

	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func AcquirePrinter() Interceptor {
	matcher := message.DefaultCatalog.Matcher()
	return func(ctx context.Context, update *tg.Update, next model.UpdateFunc) error {
		u, _ := ctx.Value("user").(*model.User)
		tag, _, _ := matcher.Match(language.Make(u.LanguageCode))
		printer := message.NewPrinter(tag)
		ctx = context.WithValue(ctx, "printer", printer)
		return next(ctx, update)
	}
}