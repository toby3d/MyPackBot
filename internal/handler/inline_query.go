package handler

import (
	"context"
	"strconv"
	"strings"

	"github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

//nolint: gochecknoglobals
var replacer = strings.NewReplacer(
	"personal:true", "",
	"personal:false", "",
)

func (h *Handler) IsInlineQuery(ctx context.Context, inline *tg.InlineQuery) (err error) {
	u, _ := ctx.Value(common.ContextUser).(*model.User)
	p, _ := ctx.Value(common.ContextPrinter).(*message.Printer)
	offset, _ := strconv.Atoi(inline.Offset)

	answer := tg.NewAnswerInlineQuery(inline.ID)
	answer.IsPersonal = !strings.Contains(inline.Query, "personal:false")
	answer.CacheTime = 1

	if inline.HasQuery() {
		inline.Query = replacer.Replace(inline.Query)
		inline.Query = strings.TrimSpace(inline.Query)
		inline.Query, _ = utils.FixEmojiTone(inline.Query)
	}

	answer.SwitchPrivateMessageText = p.Sprintf("inline__not-found_switch-text")
	answer.SwitchPrivateMessageParameter = "from_inline"

	var count int
	var stickers model.Stickers

	dlog.Ln("is_personal:", answer.IsPersonal)
	if answer.IsPersonal {
		stickers, count = h.store.GetStickersList(u, offset, 50, inline.Query)
	} else {
		stickers, count = h.store.Stickers().GetList(offset, 50, inline.Query)
	}
	dlog.Ln("count:", count)
	dlog.D(stickers)

	if count > 0 && offset+50 < count {
		answer.NextOffset = strconv.Itoa(offset + 50)
	}

	answer.SwitchPrivateMessageText = p.Sprintf("inline__found_switch-text", count)
	answer.SwitchPrivateMessageParameter = "from_inline"

	answer.Results = make([]interface{}, len(stickers))

	for i := range stickers {
		answer.Results[i] = tg.NewInlineQueryResultCachedSticker(stickers[i].ID, stickers[i].ID)
	}

	_, err = h.bot.AnswerInlineQuery(answer)

	return err
}
