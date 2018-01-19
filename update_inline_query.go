package main

import (
	"strconv"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func updateInlineQuery(inlineQuery *tg.InlineQuery) {
	fixedQuery, err := fixEmoji(inlineQuery.Query)
	if err == nil {
		inlineQuery.Query = fixedQuery
	}

	answer := &tg.AnswerInlineQueryParameters{}
	answer.InlineQueryID = inlineQuery.ID
	answer.CacheTime = 1
	answer.IsPersonal = true

	if len([]rune(inlineQuery.Query)) >= 256 {
		_, err = bot.AnswerInlineQuery(answer)
		errCheck(err)
		return
	}

	log.Ln("Let's preparing answer...")
	T, err := switchLocale(inlineQuery.From.LanguageCode)
	errCheck(err)

	log.Ln("INLINE OFFSET:", inlineQuery.Offset)
	if inlineQuery.Offset == "" {
		inlineQuery.Offset = "-1"
	}
	offset, err := strconv.Atoi(inlineQuery.Offset)
	errCheck(err)
	offset++

	stickers, packSize, err := dbGetUserStickers(
		inlineQuery.From.ID, offset, inlineQuery.Query,
	)
	errCheck(err)

	totalStickers := len(stickers)
	if totalStickers == 0 {
		if offset == 0 {
			if inlineQuery.Query != "" {
				// If search stickers by emoji return 0 results
				answer.SwitchPrivateMessageText = T(
					"button_inline_nothing", map[string]interface{}{
						"Query": inlineQuery.Query,
					},
				)
				answer.SwitchPrivateMessageParameter = cmdAddSticker
			} else {
				// If query is empty and get 0 stickers
				answer.SwitchPrivateMessageText = T("button_inline_empty")
				answer.SwitchPrivateMessageParameter = cmdAddSticker
			}
			answer.Results = nil
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
			results[i] = tg.NewInlineQueryResultCachedSticker(sticker, sticker)
		}

		answer.SwitchPrivateMessageText = T(
			"button_inline_search", packSize, map[string]interface{}{
				"Count": packSize,
			},
		)
		answer.SwitchPrivateMessageParameter = cmdHelp
		answer.Results = results
	}

	log.Ln("CacheTime:", answer.CacheTime)

	_, err = bot.AnswerInlineQuery(answer)
	errCheck(err)
}
