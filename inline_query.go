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
	if offset == -1 {
		offset++
	}

	answer := &telegram.AnswerInlineQueryParameters{
		InlineQueryID: inline.ID,
		CacheTime:     1,
		IsPersonal:    true,
	}

	stickers, emojis, err := dbGetUserStickers(inline.From.ID)
	errCheck(err)

	packSize := len(stickers)

	if inline.Query != "" {
		var buffer []string
		for i := range stickers {
			if emojis[i] != inline.Query {
				continue
			}

			buffer = append(buffer, stickers[i])
		}
		stickers = buffer
	}

	totalStickers := len(stickers)
	totalPages := totalStickers / perPage
	if totalStickers == 0 {
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
		log.Ln("LESS THAN:", offset < totalPages)
		if offset < totalPages {
			from := offset * perPage
			if offset > 0 {
				from--
			}
			to := from + perPage

			log.Ln("from:", from)
			log.Ln("to:", to)

			stickers = stickers[from:to]

			offset++
			answer.NextOffset = strconv.Itoa(offset)
		} else {
			from := offset * perPage
			if offset > 0 {
				from--
			}
			to := from
			log.Ln("MINUS:", totalStickers%perPage)
			if totalStickers%perPage != 0 {
				log.Ln("FUCK")
				to += totalStickers % perPage
			} else {
				to += perPage
			}

			log.Ln("from:", from)
			log.Ln("to:", to)
			stickers = stickers[from:to]
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
			})
		answer.SwitchPrivateMessageParameter = cmdAddSticker
		answer.Results = results
	}

	_, err = bot.AnswerInlineQuery(answer)
	errCheck(err)
}
