package events

import (
	"strconv"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (e *Events) InlineQuery(b *tg.Bot, iq *tg.InlineQuery) error {
	offset, _ := strconv.Atoi(iq.Offset)
	p := message.NewPrinter(language.English)

	answer := tg.NewAnswerInlineQuery(iq.ID)
	answer.CacheTime = 1
	answer.IsPersonal = true

	if len([]rune(iq.Query)) >= 256 {
		_, err := b.AnswerInlineQuery(answer)
		return err
	}

	answer.NextOffset = strconv.Itoa(offset + 50)
	answer.SwitchPrivateMessageText = p.Sprintf("Not found any stickers")
	answer.SwitchPrivateMessageParameter = strconv.Itoa(iq.From.ID)

	user, err := e.store.GetOrCreateUser(&models.User{
		ID:           iq.From.ID,
		LanguageCode: iq.From.LanguageCode,
		StartedAt:    time.Now().UTC().Unix(),
		AutoSaving:   true,
	})
	if err != nil {
		b.AnswerInlineQuery(answer)
		return err
	}

	if iq.Query, err = utils.FixEmojiTone(iq.Query); err != nil {
		_, err := b.AnswerInlineQuery(answer)
		return err
	}

	p = message.NewPrinter(language.Make(user.LanguageCode))

	var stickers []models.Sticker
	var count int
	if !iq.HasQuery() {
		stickers, count = e.store.GetUserStickers(user.ID, offset, 50)
	} else {
		stickers, count = e.store.GetUserStickersByQuery(iq.Query, user.ID, offset, 50)
	}

	answer.SwitchPrivateMessageText = p.Sprintf("Found %d results", count)
	for _, s := range stickers {
		s := s
		answer.Results = append(
			answer.Results,
			tg.NewInlineQueryResultCachedSticker(strconv.Itoa(user.ID)+"_"+s.ID, s.ID),
		)
	}

	_, err = b.AnswerInlineQuery(answer)
	return err
}
