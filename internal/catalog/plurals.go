package catalog

import (
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func RegisterPlurals() (err error) {
	if err = message.Set(language.English, "\u200d Found %d result(s)", plural.Selectf(1, "%d",
		"=0", "🤷 Nothing found",
		"=1", "\u200d Found %d result",
		"other", "\u200d Found %[1]d results",
	)); err != nil {
		return err
	}

	return message.Set(language.Russian, "\u200d Found %d result(s)", plural.Selectf(1, "%d",
		"=0", "🤷 Ничего не найдено",
		"=1", "\u200d Найден %[1]d результат",
		"<5", "\u200d Найдено %[1]d результата",
		"other", "\u200d Найдено %[1]d результатов",
	))
}
