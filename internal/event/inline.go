package event

import (
	"strconv"
	"strings"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/middleware"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

func (event *Event) InlineQuery(i *tg.InlineQuery) (err error) {
	answer := tg.NewAnswerInlineQuery(i.ID)
	answer.IsPersonal = !strings.Contains(i.Query, "personal:false")
	answer.CacheTime = 1
	if i.HasQuery() {
		i.Query = strings.Trim(i.Query, "personal:false")
		i.Query = strings.Trim(i.Query, "personal:true")
		i.Query, _ = utils.FixEmojiTone(i.Query)
		i.Query = strings.TrimSpace(i.Query)
	}

	u, err := middleware.GetUser(event.store, i.From, time.Now().UTC().Unix())
	if err != nil {
		_, err = event.bot.AnswerInlineQuery(answer)
		return err
	}

	p := middleware.GetPrinter(u.LanguageCode)
	answer.SwitchPrivateMessageText = p.Sprintf("inline__not-found_switch-text")
	answer.SwitchPrivateMessageParameter = "from_inline"

	offset, _ := strconv.Atoi(i.Offset)
	stickers, count := event.store.GetStickersList(u, offset, 50, i.Query)
	if !answer.IsPersonal {
		stickers, count = event.store.Stickers().GetList(offset, 50, i.Query)
	}

	if count != 0 && offset+50 < count {
		answer.NextOffset = strconv.Itoa(offset + 50)
	}

	answer.SwitchPrivateMessageText = p.Sprintf("inline__found_switch-text", count)
	answer.SwitchPrivateMessageParameter = "from_inline"

	answer.Results = make([]interface{}, len(stickers), len(stickers))
	for i := range stickers {
		answer.Results[i] = tg.NewInlineQueryResultCachedSticker(stickers[i].ID, stickers[i].ID)
	}

	_, err = event.bot.AnswerInlineQuery(answer)
	return err
}
