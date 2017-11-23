package main

import (
	"strconv"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

func inlineQuery(inline *telegram.InlineQuery) {
	log.Ln("Let's preparing answer...")
	T, err := switchLocale(inline.From.LanguageCode)
	errCheck(err)

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
		answer.SwitchPrivateMessageText = T("button_inline_empty")
		answer.SwitchPrivateMessageParameter = cmdAddSticker
	case offset <= 0 && len(stickers) == 0 && inline.Query != "":
		// If search stickers by emoji return 0 results
		answer.SwitchPrivateMessageText = T(
			"button_inline_nothing",
			map[string]interface{}{"Query": inline.Query},
		)
		answer.SwitchPrivateMessageParameter = cmdAddSticker
	case offset >= 0 && len(stickers) == 50,
		offset >= 0 && len(stickers) < 50:
		offset++
		answer.NextOffset = strconv.Itoa(offset)

		var results = make([]interface{}, len(stickers))
		for i, sticker := range stickers {
			results[i] = telegram.NewInlineQueryResultCachedSticker(
				sticker, // resultID
				sticker, // fileID
			)
		}

		answer.SwitchPrivateMessageText = T("button_inline_add")
		answer.SwitchPrivateMessageParameter = cmdAddSticker
		answer.Results = results
	}

	_, err = bot.AnswerInlineQuery(&answer)
	errCheck(err)
}
