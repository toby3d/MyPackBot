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

const defaultLimit int = 50

func (h *Handler) IsInlineQuery(ctx *model.Context) (err error) {
	answer := tg.NewAnswerInline(ctx.Request.InlineQuery.ID)
	answer.CacheTime = 1 // NOTE(toby3d): add setting for change this
	filter := getFilter(ctx.Request.InlineQuery)

	if answer.IsPersonal = !strings.Contains(ctx.Request.InlineQuery.Query, "personal:false"); answer.IsPersonal {
		filter.UserID = ctx.User.ID
	}

	results, count, _ := h.store.GetList(filter.Offset, filter.Limit, filter)

	if filter.Offset+filter.Limit < count {
		answer.NextOffset = strconv.Itoa(filter.Offset + filter.Limit)
	}

	for i := range results {
		switch results[i].GetType() {
		case tg.TypeSticker:
			answer.Results = append(answer.Results, tg.NewInlineQueryResultCachedSticker(
				tg.TypeSticker+common.DataSeparator+results[i].GetID(), results[i].GetFileID(),
			))
		case tg.TypePhoto:
			answer.Results = append(answer.Results, tg.NewInlineQueryResultCachedPhoto(
				tg.TypePhoto+common.DataSeparator+results[i].GetID(), results[i].GetFileID(),
			))
		}
	}

	p := ctx.Get("printer").(*message.Printer)
	answer.SwitchPrivateMessageParameter = "inline"
	answer.SwitchPrivateMessageText = p.Sprintf("🕵 Found %d result(s)", count)
	_, err = ctx.AnswerInlineQuery(answer)

	return err
}

func getFilter(iq *tg.InlineQuery) *store.Filter {
	f := new(store.Filter)
	f.Limit = defaultLimit
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
		/*
			case strings.HasPrefix(field, "animated:"):
				isAnimated, _ := strconv.ParseBool(strings.TrimPrefix(field, "animated:"))
				f.IsAnimated = &isAnimated
			case strings.HasPrefix(field, "set:"):
				f.SetName = strings.TrimPrefix(field, "set:")
		*/
		case strings.HasPrefix(field, "photos:"), strings.HasPrefix(field, "stickers:"):
		default:
			continue
		}

		f.Query = f.Query[:i] + strings.TrimPrefix(f.Query[i:], field)
	}

	f.Query = strings.ReplaceAll(f.Query, " ", "")

	return f
}
