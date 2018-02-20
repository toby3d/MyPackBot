package updates

import (
	"strconv"

	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/helpers"
	"github.com/toby3d/MyPackBot/internal/i18n"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// InlineQuery checks InlineQuery updates for answer with personal results
func InlineQuery(inlineQuery *tg.InlineQuery) {
	fixedQuery, err := helpers.FixEmoji(inlineQuery.Query)
	if err == nil {
		inlineQuery.Query = fixedQuery
	}

	answer := &tg.AnswerInlineQueryParameters{}
	answer.InlineQueryID = inlineQuery.ID
	answer.CacheTime = 1
	answer.IsPersonal = true

	if len([]rune(inlineQuery.Query)) >= 256 {
		_, err = bot.Bot.AnswerInlineQuery(answer)
		errors.Check(err)
		return
	}

	log.Ln("Let's preparing answer...")
	T, err := i18n.SwitchTo(inlineQuery.From.LanguageCode)
	errors.Check(err)

	log.Ln("INLINE OFFSET:", inlineQuery.Offset)
	if inlineQuery.Offset == "" {
		inlineQuery.Offset = "-1"
	}
	offset, err := strconv.Atoi(inlineQuery.Offset)
	errors.Check(err)
	offset++

	stickers, packSize, err := db.UserStickers(
		inlineQuery.From.ID, offset, inlineQuery.Query,
	)
	errors.Check(err)

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
				answer.SwitchPrivateMessageParameter = models.CommandAddSticker
			} else {
				// If query is empty and get 0 stickers
				answer.SwitchPrivateMessageText = T("button_inline_empty")
				answer.SwitchPrivateMessageParameter = models.CommandAddSticker
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
		answer.SwitchPrivateMessageParameter = models.CommandHelp
		answer.Results = results
	}

	log.Ln("CacheTime:", answer.CacheTime)

	_, err = bot.Bot.AnswerInlineQuery(answer)
	errors.Check(err)
}
