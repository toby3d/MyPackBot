package main

import (
	"fmt"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

func inlineQuery(inline *telegram.InlineQuery) {
	answer := telegram.AnswerInlineQueryParameters{
		InlineQueryID:                 inline.ID,
		CacheTime:                     1,
		IsPersonal:                    true,
		SwitchPrivateMessageText:      "No stickers, let's add some one!",
		SwitchPrivateMessageParameter: "addSticker",
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
			)
		}

		answer.SwitchPrivateMessageText = "Add one more sticker!"
		answer.SwitchPrivateMessageParameter = "addSticker"
		answer.Results = results
	}

	_, err = bot.AnswerInlineQuery(&answer)
	errCheck(err)
}
