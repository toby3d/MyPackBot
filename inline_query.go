package main

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func inlineQuery(inline *telegram.InlineQuery) {
	T, err := i18n.Tfunc(inline.From.LanguageCode)
	if err != nil {
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}
	answer := telegram.AnswerInlineQueryParameters{
		InlineQueryID:                 inline.ID,
		CacheTime:                     1,
		IsPersonal:                    true,
	}

	stickers, err := dbGetUserStickers(inline.From.ID, inline.Query)
	errCheck(err)

	if len(stickers) > 0 {
		var results []interface{}
		for i := range stickers {
			results = append(
				results,
				telegram.NewInlineQueryResultCachedSticker(
					fmt.Sprint("sticker", stickers[i]), // resultID
					stickers[i],                        // fileID
				),
		answer.SwitchPrivateMessageText = T("inline_empty")
		answer.SwitchPrivateMessageParameter = "add"
		answer.SwitchPrivateMessageText = T(
			"inline_nothing",
			map[string]interface{}{
				"Query": inline.Query,
			},
		)
		answer.SwitchPrivateMessageParameter = "help"
			)
		}

		answer.Results = results
	}

	_, err = bot.AnswerInlineQuery(&answer)
	errCheck(err)
}
