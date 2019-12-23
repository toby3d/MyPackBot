package handler

import (
	"strconv"
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) IsInlineQuery(ctx *model.Context) (err error) {
	p := ctx.Get("printer").(*message.Printer)
	answer := tg.NewAnswerInlineQuery(ctx.Request.InlineQuery.ID)
	answer.IsPersonal = !strings.Contains(ctx.Request.InlineQuery.Query, "personal:false")
	answer.CacheTime = 1 // NOTE(toby3d): add setting for change this

	if ctx.Request.InlineQuery.HasQuery() {
		ctx.Request.InlineQuery.Query = strings.TrimSpace(ctx.Request.InlineQuery.Query)
		ctx.Request.InlineQuery.Query, _ = utils.FixEmojiTone(ctx.Request.InlineQuery.Query)
	}

	offset, _ := strconv.Atoi(ctx.Request.InlineQuery.Offset)
	items, count := h.store.GetList(ctx.User.ID, offset, 50, ctx.Request.InlineQuery.Query)

	if count > offset+50 {
		answer.NextOffset = strconv.Itoa(offset + 50)
	}

	for i := range items {
		switch item := items[i].(type) {
		case *model.Sticker:
			answer.Results = append(answer.Results, tg.NewInlineQueryResultCachedSticker(
				tg.TypeSticker+common.DataSeparator+strconv.FormatUint(item.ID, 10), item.FileID,
			))
		case *model.Photo:
			answer.Results = append(answer.Results, tg.NewInlineQueryResultCachedPhoto(
				tg.TypePhoto+common.DataSeparator+strconv.FormatUint(item.ID, 10), item.FileID,
			))
		}
	}

	p := ctx.Get("printer").(*message.Printer)
	answer.SwitchPrivateMessageParameter = "inline"
	answer.SwitchPrivateMessageText = p.Sprintf("\u200d Found %d result(s)", count)
	_, err = ctx.AnswerInlineQuery(answer)

	return err
}
