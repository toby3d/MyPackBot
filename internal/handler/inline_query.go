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
	offset, _ := strconv.Atoi(ctx.InlineQuery.Offset)
	answer := tg.NewAnswerInlineQuery(ctx.InlineQuery.ID)
	answer.IsPersonal = !strings.Contains(ctx.InlineQuery.Query, "personal:false")
	answer.CacheTime = 1
	answer.SwitchPrivateMessageText = ctx.T().Sprintf("inline__not-found_switch-text")
	answer.SwitchPrivateMessageParameter = "from_inline"

	if ctx.InlineQuery.HasQuery() {
		ctx.InlineQuery.Query = replacer.Replace(ctx.InlineQuery.Query)
		ctx.InlineQuery.Query = strings.TrimSpace(ctx.InlineQuery.Query)
		ctx.InlineQuery.Query, _ = utils.FixEmojiTone(ctx.InlineQuery.Query)
	}

	var stickers model.Stickers

	count := 0
	if answer.IsPersonal {
		stickers, count = h.store.GetStickersList(ctx.User, offset, 50, ctx.InlineQuery.Query)
	} else {
		stickers, count = h.store.Stickers().GetList(offset, 50, ctx.InlineQuery.Query)
	}

	if count > 0 && offset+50 < count {
		answer.NextOffset = strconv.Itoa(offset + 50)
	}

	answer.SwitchPrivateMessageText = ctx.T().Sprintf("inline__found_switch-text", count)
	answer.SwitchPrivateMessageParameter = "from_inline"
	answer.Results = make([]interface{}, len(stickers))

	for i := range stickers {
		answer.Results[i] = tg.NewInlineQueryResultCachedSticker(stickers[i].ID, stickers[i].ID)
	}

	_, err = ctx.AnswerInlineQuery(answer)

	return err
}
