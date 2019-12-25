package handler

import (
	"strconv"
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/store"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) IsInlineQuery(ctx *model.Context) (err error) {
	answer := tg.NewAnswerInlineQuery(ctx.Request.InlineQuery.ID)
	answer.CacheTime = 1 // NOTE(toby3d): add setting for change this
	answer.IsPersonal = !strings.Contains(ctx.Request.InlineQuery.Query, "personal:false")
	filter := getFilter(ctx.Request.InlineQuery)
	items, count := h.store.GetList(ctx.User.ID, filter)

	if filter.Offset+50 < count {
		answer.NextOffset = strconv.Itoa(filter.Offset + 50)
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

func getFilter(iq *tg.InlineQuery) *store.Filter {
	f := new(store.Filter)
	f.Limit = 50
	f.Offset, _ = strconv.Atoi(iq.Offset)

	if !strings.Contains(iq.Query, "photos:false") {
		f.AllowedTypes = append(f.AllowedTypes, tg.TypePhoto)
	}

	if !strings.Contains(iq.Query, "stickers:false") {
		f.AllowedTypes = append(f.AllowedTypes, tg.TypeSticker)
	}

	if !iq.HasQuery() {
		return f
	}

	f.Query, _ = utils.FixEmojiTone(strings.TrimSpace(iq.Query))
	for _, field := range strings.Fields(f.Query) {
		i := strings.Index(f.Query, field)

		switch {
		case strings.HasPrefix(field, "offset:"):
			f.Offset, _ = strconv.Atoi(strings.TrimPrefix(field, "offset:"))
		case strings.HasPrefix(field, "animated:"):
			isAnimated, _ := strconv.ParseBool(strings.TrimPrefix(field, "animated:"))
			f.IsAnimated = &isAnimated
		case strings.HasPrefix(field, "set:"):
			f.SetName = strings.TrimPrefix(field, "set:")
		case strings.HasPrefix(field, "photos:"), strings.HasPrefix(field, "stickers:"):
		default:
			continue
		}

		f.Query = f.Query[:i] + strings.TrimPrefix(f.Query[i:], field)
	}

	f.Query = strings.ReplaceAll(f.Query, " ", "")

	return f
}
