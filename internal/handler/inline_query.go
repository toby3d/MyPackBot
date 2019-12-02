package handler

import (
	"strconv"
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

func (h *Handler) IsInlineQuery(ctx *model.Context) (err error) {
	offset, _ := strconv.Atoi(ctx.Request.InlineQuery.Offset)
	answer := tg.NewAnswerInlineQuery(ctx.Request.InlineQuery.ID)
	answer.IsPersonal = !strings.Contains(ctx.Request.InlineQuery.Query, "personal:false")
	answer.SwitchPrivateMessageText = "inline__not-found_switch-text"
	answer.SwitchPrivateMessageParameter = "from_inline"
	answer.CacheTime = 1

	if ctx.Request.InlineQuery.HasQuery() {
		ctx.Request.InlineQuery.Query = strings.TrimSpace(ctx.Request.InlineQuery.Query)
		ctx.Request.InlineQuery.Query, _ = utils.FixEmojiTone(ctx.Request.InlineQuery.Query)
	}

	items, count := h.store.GetList(ctx.User.ID, offset, 50, ctx.Request.InlineQuery.Query)
	answer.SwitchPrivateMessageText = "inline__found_switch-text " + strconv.Itoa(count)

	if count > offset+50 {
		answer.NextOffset = strconv.Itoa(offset + 50)
	}

	results := make([]interface{}, 0, 50)
	for i := range items {
		switch item := items[i].(type) {
		case *model.Sticker:
			results = append(results, tg.NewInlineQueryResultCachedSticker(
				"sticker@"+strconv.FormatUint(item.ID, 10), item.FileID,
			))
		case *model.Photo:
			results = append(results, tg.NewInlineQueryResultCachedPhoto(
				"photo@"+strconv.FormatUint(item.ID, 10), item.FileID,
			))
		}
	}
	answer.Results = results

	_, err = ctx.AnswerInlineQuery(answer)

	return err
}
