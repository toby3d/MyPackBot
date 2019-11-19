package handler

import (
	"strconv"
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

//nolint: gochecknoglobals
var replacer = strings.NewReplacer(
	"personal:true", "",
	"personal:false", "",
)

func (h *Handler) IsInlineQuery(ctx *model.Context) (err error) {
	var stickers model.Stickers

	answer := tg.NewAnswerInlineQuery(ctx.InlineQuery.ID)
	answer.IsPersonal = !strings.Contains(ctx.InlineQuery.Query, "personal:false")
	answer.SwitchPrivateMessageText = ctx.T().Sprintf("inline__not-found_switch-text")
	answer.SwitchPrivateMessageParameter = "from_inline"
	answer.CacheTime = 1
	offset, _ := strconv.Atoi(ctx.InlineQuery.Offset)
	count := 0

	if ctx.InlineQuery.HasQuery() {
		ctx.InlineQuery.Query = replacer.Replace(ctx.InlineQuery.Query)
		ctx.InlineQuery.Query = strings.TrimSpace(ctx.InlineQuery.Query)
		ctx.InlineQuery.Query, _ = utils.FixEmojiTone(ctx.InlineQuery.Query)
	}

	if answer.IsPersonal {
		stickers, count = h.store.GetStickersList(ctx.User, offset, 50, ctx.InlineQuery.Query)
	} else {
		stickers, count = h.stickersStore.GetList(offset, 50, ctx.InlineQuery.Query)
	}

	if count > offset+50 {
		answer.NextOffset = strconv.Itoa(offset + 50)
	}

	answer.SwitchPrivateMessageText = ctx.T().Sprintf("inline__found_switch-text", count)
	answer.SwitchPrivateMessageParameter = "from_inline"
	answer.Results = make([]interface{}, len(stickers))

	for i := range stickers {
		answer.Results[i] = tg.NewInlineQueryResultCachedSticker(stickers[i].ID, stickers[i].ID)
	}

	_, err = ctx.AnswerInlineQuery(answer)

	return ctx.Error(err)
}
