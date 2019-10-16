package handler

import (
	"context"
	"strconv"
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (h *Handler) isInlineQuery(ctx context.Context, inline *tg.InlineQuery) (err error) {
	u, _ := ctx.Value("user").(*model.User)
	p, _ := ctx.Value("printer").(*message.Printer)

	answer := tg.NewAnswerInlineQuery(inline.ID)
	answer.IsPersonal = !strings.Contains(inline.Query, "personal:false")
	answer.CacheTime = 1
	if inline.HasQuery() {
		inline.Query = strings.Trim(inline.Query, "personal:true")
		inline.Query = strings.Trim(inline.Query, "personal:false")
		inline.Query = strings.TrimSpace(inline.Query)
		inline.Query, _ = utils.FixEmojiTone(inline.Query)
	}

	answer.SwitchPrivateMessageText = p.Sprintf("inline__not-found_switch-text")
	answer.SwitchPrivateMessageParameter = "from_inline"

	offset, _ := strconv.Atoi(inline.Offset)
	stickers, count := h.store.GetStickersList(u, offset, 50, inline.Query)
	if !answer.IsPersonal {
		stickers, count = h.store.Stickers().GetList(offset, 50, inline.Query)
	}

	if count > 0 && offset+50 < count {
		answer.NextOffset = strconv.Itoa(offset + 50)
	}

	answer.SwitchPrivateMessageText = p.Sprintf("inline__found_switch-text", count)
	answer.SwitchPrivateMessageParameter = "from_inline"

	answer.Results = make([]interface{}, len(stickers), len(stickers))
	for i := range stickers {
		answer.Results[i] = tg.NewInlineQueryResultCachedSticker(stickers[i].ID, stickers[i].ID)
	}

	_, err = h.bot.AnswerInlineQuery(answer)
	return err
}
