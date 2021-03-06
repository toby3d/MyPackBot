package updates

import (
	"strconv"

	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// InlineQuery checks InlineQuery updates for answer with personal results
func InlineQuery(inlineQuery *tg.InlineQuery) {
	fixedQuery, err := utils.FixEmoji(inlineQuery.Query)
	if err == nil {
		inlineQuery.Query = fixedQuery
	}

	answer := new(tg.AnswerInlineQueryParameters)
	answer.InlineQueryID = inlineQuery.ID
	answer.CacheTime = 1
	answer.IsPersonal = true

	if len([]rune(inlineQuery.Query)) >= 256 {
		_, err = bot.Bot.AnswerInlineQuery(answer)
		errors.Check(err)
		return
	}

	log.Ln("Let's preparing answer...")
	t, err := i18n.SwitchTo(inlineQuery.From.LanguageCode)
	errors.Check(err)

	log.Ln("INLINE OFFSET:", inlineQuery.Offset)
	if inlineQuery.Offset == "" {
		inlineQuery.Offset = "-1"
	}
	offset, err := strconv.Atoi(inlineQuery.Offset)
	errors.Check(err)
	offset++

	stickers, err := db.DB.GetUserStickers(
		inlineQuery.From,
		&tg.InlineQuery{
			Offset: strconv.Itoa(offset),
			Query:  inlineQuery.Query,
		},
	)
	errors.Check(err)

	if len(stickers) == 0 {
		if offset == 0 && inlineQuery.Query != "" {
			// If search stickers by emoji return 0 results
			answer.SwitchPrivateMessageText = t(
				"button_inline_nothing", map[string]interface{}{
					"Query": inlineQuery.Query,
				},
			)

			answer.SwitchPrivateMessageParameter = tg.CommandHelp
		}

		answer.Results = nil
	} else {
		log.Ln("STICKERS FROM REQUEST:", len(stickers))
		if len(stickers) == 50 {
			answer.NextOffset = strconv.Itoa(offset)
			log.Ln("NEXT OFFSET:", answer.NextOffset)
		}

		var results = make([]interface{}, len(stickers))
		for i, sticker := range stickers {
			results[i] = tg.NewInlineQueryResultCachedSticker(sticker, sticker)
		}

		answer.Results = results
	}

	log.Ln("CacheTime:", answer.CacheTime)

	_, err = bot.Bot.AnswerInlineQuery(answer)
	errors.Check(err)
}
