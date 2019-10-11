package event

import (
	"strings"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/middleware"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
)

func (event *Event) CallbackQuery(c *tg.CallbackQuery) error {
	u, err := middleware.GetUser(event.store, c.From, time.Now().UTC().Unix())
	if err != nil {
		return err
	}

	parts := strings.Split(c.Data, ":")
	switch parts[0] {
	case "add":
		switch parts[1] {
		case "single":
			return event.CallbackAddSingleSticker(u, c)
		case "set":
			return event.CallbackAddStickerSet(u, c)
		}
	case "remove":
		switch parts[1] {
		case "single":
			return event.CallbackRemoveSingleSticker(u, c)
		case "set":
			return event.CallbackRemoveStickerSet(u, c)
		}
	}
	return nil
}

func (event *Event) CallbackAddSingleSticker(u *model.User, c *tg.CallbackQuery) (err error) {
	answer := tg.NewAnswerCallbackQuery(c.ID)
	p := middleware.GetPrinter(u.LanguageCode)
	answer.Text = p.Sprintf("callback__text_add-single", c.Message.ReplyToMessage.Sticker.SetName)
	if err = event.store.AddSticker(u, &model.Sticker{
		CreatedAt:  c.Message.ReplyToMessage.Date,
		Emoji:      c.Message.ReplyToMessage.Sticker.Emoji,
		ID:         c.Message.ReplyToMessage.Sticker.FileID,
		IsAnimated: c.Message.ReplyToMessage.Sticker.IsAnimated,
		SetName:    c.Message.ReplyToMessage.Sticker.SetName,
	}); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = event.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-single"), "remove:single"),
	))
	if len(c.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], c.Message.ReplyMarkup.InlineKeyboard[0][1])
	}

	if _, err = event.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          c.Message.Chat.ID,
		InlineMessageID: c.InlineMessageID,
		MessageID:       c.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = event.bot.AnswerCallbackQuery(answer)
	return err
}

func (event *Event) CallbackAddStickerSet(u *model.User, c *tg.CallbackQuery) (err error) {
	answer := tg.NewAnswerCallbackQuery(c.ID)
	set, err := event.bot.GetStickerSet(c.Message.ReplyToMessage.Sticker.SetName)
	if err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = event.bot.AnswerCallbackQuery(answer)
		return err
	}
	for i := range set.Stickers {
		s, err := event.store.Stickers().GetOrCreate(&model.Sticker{
			CreatedAt:  time.Now().UTC().Unix(),
			Emoji:      set.Stickers[i].Emoji,
			ID:         set.Stickers[i].FileID,
			IsAnimated: set.Stickers[i].IsAnimated,
			SetName:    set.Name,
		})
		if err != nil {
			answer.Text = "🐞 " + err.Error()
			_, err = event.bot.AnswerCallbackQuery(answer)
			return err
		}
		if err = event.store.AddSticker(u, s); err != nil {
			answer.Text = "🐞 " + err.Error()
			_, err = event.bot.AnswerCallbackQuery(answer)
			return err
		}
	}

	p := middleware.GetPrinter(u.LanguageCode)
	answer.Text = p.Sprintf("callback__text_add-set", c.Message.ReplyToMessage.Sticker.SetName)
	if err = event.store.AddStickersSet(u, c.Message.ReplyToMessage.Sticker.SetName); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = event.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-single"), "remove:single"),
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-set"), "remove:set"),
	))
	if _, err = event.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          c.Message.Chat.ID,
		InlineMessageID: c.InlineMessageID,
		MessageID:       c.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = event.bot.AnswerCallbackQuery(answer)
	return err
}

func (event *Event) CallbackRemoveSingleSticker(u *model.User, c *tg.CallbackQuery) (err error) {
	answer := tg.NewAnswerCallbackQuery(c.ID)
	p := middleware.GetPrinter(u.LanguageCode)
	answer.Text = p.Sprintf("callback__text_remove-single")
	if err = event.store.RemoveSticker(u, &model.Sticker{
		CreatedAt:  c.Message.ReplyToMessage.Date,
		Emoji:      c.Message.ReplyToMessage.Sticker.Emoji,
		ID:         c.Message.ReplyToMessage.Sticker.FileID,
		IsAnimated: c.Message.ReplyToMessage.Sticker.IsAnimated,
		SetName:    c.Message.ReplyToMessage.Sticker.SetName,
	}); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = event.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-single"), "add:single"),
	))
	if len(c.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], c.Message.ReplyMarkup.InlineKeyboard[0][1])
	}

	if _, err = event.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          c.Message.Chat.ID,
		InlineMessageID: c.InlineMessageID,
		MessageID:       c.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = event.bot.AnswerCallbackQuery(answer)
	return err
}

func (event *Event) CallbackRemoveStickerSet(u *model.User, c *tg.CallbackQuery) (err error) {
	answer := tg.NewAnswerCallbackQuery(c.ID)
	p := middleware.GetPrinter(u.LanguageCode)
	answer.Text = p.Sprintf("callback__text_remove-set", c.Message.ReplyToMessage.Sticker.SetName)
	if err = event.store.RemoveStickersSet(u, c.Message.ReplyToMessage.Sticker.SetName); err != nil {
		answer.Text = "🐞 " + err.Error()
		_, err = event.bot.AnswerCallbackQuery(answer)
		return err
	}

	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-single"), "add:single"),
	))
	if len(c.Message.ReplyMarkup.InlineKeyboard[0]) == 2 {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-set"), "add:set"),
		)
	}
	if _, err = event.bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
		ChatID:          c.Message.Chat.ID,
		InlineMessageID: c.InlineMessageID,
		MessageID:       c.Message.ID,
		ReplyMarkup:     markup,
	}); err != nil {
		return
	}

	_, err = event.bot.AnswerCallbackQuery(answer)
	return err
}
