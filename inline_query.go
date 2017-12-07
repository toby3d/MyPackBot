package main

import (
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

const perPage = 50

var r = strings.NewReplacer(
	"🏻", "",
	"🏼", "",
	"🏽", "",
	"🏾", "",
	"🏿", "",
)

func inlineQuery(inline *telegram.InlineQuery) {
	inline.Query = r.Replace(inline.Query)

	log.Ln("Let's preparing answer...")
	T, err := switchLocale(inline.From.LanguageCode)
	errCheck(err)

	log.Ln("INLINE OFFSET:", inline.Offset)
	if inline.Offset == "" {
		inline.Offset = "-1"
	}
	offset, err := strconv.Atoi(inline.Offset)
	errCheck(err)
	offset++

	log.Ln("CURRENT OFFSET:", inline.Offset)
	answer := &telegram.AnswerInlineQueryParameters{}
	answer.InlineQueryID = inline.ID
	answer.CacheTime = 1
	answer.IsPersonal = true

	stickers, packSize, err := dbGetUserStickers(inline.From.ID, offset, inline.Query)
	errCheck(err)

	totalStickers := len(stickers)
	if totalStickers == 0 {
		if offset == 0 {
			if inline.Query != "" {
				// If search stickers by emoji return 0 results
				answer.SwitchPrivateMessageText = T(
					"button_inline_nothing",
					map[string]interface{}{"Query": inline.Query},
				)
				answer.SwitchPrivateMessageParameter = cmdAddSticker
			} else {
				// If query is empty and get 0 stickers
				answer.SwitchPrivateMessageText = T("button_inline_empty")
				answer.SwitchPrivateMessageParameter = cmdAddSticker
			}
		} else {
			return
		}
	} else {
		log.Ln("STICKERS FROM REQUEST:", totalStickers)
		if totalStickers > 50 {
			answer.NextOffset = strconv.Itoa(offset)
			log.Ln("NEXT OFFSET:", answer.NextOffset)

			stickers = stickers[:totalStickers-1]
		}

		log.Ln("Stickers after checks:", len(stickers))

		var results = make([]interface{}, len(stickers))
		for i, sticker := range stickers {
			results[i] = telegram.NewInlineQueryResultCachedSticker(
				sticker, sticker,
			)
		}

		answer.SwitchPrivateMessageText = T(
			"button_inline_add",
			packSize,
			map[string]interface{}{
				"Count": packSize,
			},
		)
		answer.SwitchPrivateMessageParameter = cmdAddSticker
		answer.Results = results
	}

	log.Ln("CacheTime:", answer.CacheTime)

	_, err := bot.AnswerInlineQuery(answer)
	if err != nil {
		log.Ln(err.Error())
	}
}
