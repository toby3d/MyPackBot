package main

import (
	"fmt"
	"strconv"

	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func inlineQuery(inline *telegram.InlineQuery) {
	log.Ln("[inlineQuery] Let's preparing answer...")
	T, err := i18n.Tfunc(inline.From.LanguageCode)
	if err != nil {
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	if inline.Offset == "" {
		inline.Offset = "-1"
	}
	offset, err := strconv.Atoi(inline.Offset)
	errCheck(err)
	offset++

	answer := telegram.AnswerInlineQueryParameters{
		InlineQueryID: inline.ID,
		CacheTime:     1,
		IsPersonal:    true,
	}

	stickers, err := dbGetUserStickers(inline.From.ID, offset, inline.Query)
	errCheck(err)

	switch {
	case offset <= 0 && len(stickers) == 0 && inline.Query == "":
		// If query is empty and get 0 stickers
		answer.SwitchPrivateMessageText = T("inline_empty")
		answer.SwitchPrivateMessageParameter = "add"
	case offset <= 0 && len(stickers) == 0 && inline.Query != "":
		// If search stickers by emoji return 0 results
		answer.SwitchPrivateMessageText = T(
			"inline_nothing",
			map[string]interface{}{
				"Query": inline.Query,
			},
		)
		answer.SwitchPrivateMessageParameter = "help"
	case offset >= 0 && len(stickers) == 50,
		offset >= 0 && len(stickers) < 50:
		offset++
		answer.NextOffset = strconv.Itoa(offset)

		var results = make([]interface{}, len(stickers))
		for i, sticker := range stickers {
			results[i] = telegram.NewInlineQueryResultCachedSticker(
				fmt.Sprint("sticker", sticker), // resultID
				sticker, // fileID

			)
		}

		answer.Results = results
	}

	_, err = bot.AnswerInlineQuery(&answer)
	errCheck(err)
}
