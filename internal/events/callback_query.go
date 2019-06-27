package events

import (
	"strings"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message"
)

func (e *Events) CallbackQuery(b *tg.Bot, cc *tg.CallbackQuery) error {
	p := message.NewPrinter(language.English)
	answer := tg.NewAnswerCallbackQuery(cc.ID)
	answer.ShowAlert = false

	user, err := e.store.GetOrCreateUser(&models.User{
		ID:           cc.From.ID,
		LanguageCode: cc.From.LanguageCode,
		StartedAt:    time.Now().UTC().Unix(),
		AutoSaving:   true,
	})
	if err := e.store.CreateUser(user); err != nil {
		answer.Text = "🤷🏻‍♂️ " + p.Sprintf("Unknown error")
		_, err = b.AnswerCallbackQuery(answer)
		return err
	}
	p = message.NewPrinter(language.Make(user.LanguageCode))

	data := strings.Split(cc.Data, ":")
	switch data[0] {
	// TODO: case "remove_sticker":
	// TODO: case "remove_set":
	case "add_set":
		set := e.store.GetSet(cc.Message.Sticker.SetName)
		if set == nil {
			tgSet, err := b.GetStickerSet(cc.Message.Sticker.SetName)
			if err != nil {
				_, err = b.AnswerCallbackQuery(answer)
				return err
			}

			if set, err = e.store.GetOrCreateSet(&models.Set{
				Name:  tgSet.Name,
				Title: tgSet.Title,
			}); err != nil {
				_, err = b.AnswerCallbackQuery(answer)
				return err
			}

			go func() {
				for _, tgSticker := range tgSet.Stickers {
					sticker, err := e.store.GetOrCreateSticker(&models.Sticker{
						Model: models.Model{
							ID:      tgSticker.FileID,
							SavedAt: time.Now().UTC().Unix(),
						},
						Emoji:   tgSticker.Emoji,
						SetName: set.Name,
					})
					if err != nil {
						continue
					}

					if err = e.store.AddSticker(user.ID, sticker.ID); err != nil {
						continue
					}
				}
			}()
		}
		go func() {
			if _, err = b.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
				ChatID:          cc.Message.Chat.ID,
				InlineMessageID: cc.InlineMessageID,
				MessageID:       cc.Message.ID,
				ReplyMarkup: tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButton(
						"🔥 "+p.Sprintf("Remove %s set", set.Title), "remove_set",
					),
				)),
			}); err != nil {
				return
			}
		}()

		answer.Text = "👍🏻 " + p.Sprintf("%s set has been added!", set.Title)
	case "add_sticker":
		sticker, err := e.store.GetOrCreateSticker(&models.Sticker{
			Model: models.Model{
				ID:      cc.Message.Sticker.FileID,
				SavedAt: cc.Message.Time().UTC().Unix(),
			},
			Emoji:   cc.Message.Sticker.Emoji,
			SetName: cc.Message.Sticker.SetName,
		})
		if err != nil {
			return err
		}
		if err = e.store.AddSticker(user.ID, sticker.ID); err != nil {
			_, err = b.AnswerCallbackQuery(answer)
			return err
		}

		set := e.store.GetSet(cc.Message.Sticker.SetName)
		if set == nil {
			tgSet, err := b.GetStickerSet(cc.Message.Sticker.SetName)
			if err != nil {
				_, err = b.AnswerCallbackQuery(answer)
				return err
			}

			if set, err = e.store.GetOrCreateSet(&models.Set{
				Name:  tgSet.Name,
				Title: tgSet.Title,
			}); err != nil {
				_, err = b.AnswerCallbackQuery(answer)
				return err
			}

			go func() {
				for _, tgSticker := range tgSet.Stickers {
					sticker, err := e.store.GetOrCreateSticker(&models.Sticker{
						Model: models.Model{
							ID:      tgSticker.FileID,
							SavedAt: time.Now().UTC().Unix(),
						},
						Emoji:   tgSticker.Emoji,
						SetName: set.Name,
					})
					if err != nil {
						continue
					}

					if err = e.store.AddSticker(user.ID, sticker.ID); err != nil {
						continue
					}
				}
			}()
		}
		go func() {
			if _, err = b.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
				ChatID:          cc.Message.Chat.ID,
				InlineMessageID: cc.InlineMessageID,
				MessageID:       cc.Message.ID,
				ReplyMarkup: tg.NewInlineKeyboardMarkup(
					tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButton(
						"🔥 "+p.Sprintf("Remove this sticker"), "remove_sticker",
					)),
					tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButton(
						"💕 "+p.Sprintf("Add %s set", set.Title), "add_set",
					)),
				),
			}); err != nil {
				return
			}
		}()

		answer.Text = "👍🏻 " + p.Sprintf("%s sticker has been added!", set.Title)
	case "goto_settings":
		btnLanguage := tg.NewInlineKeyboardButton(
			p.Sprintf("language_flag")+" "+p.Sprintf("Language"), "goto_languages",
		)

		btnAutosaving := tg.NewInlineKeyboardButton("☑️ "+p.Sprintf("Autosaving"), "goto_autosaving")
		if user.AutoSaving {
			btnAutosaving.Text = "✅ " + p.Sprintf("Autosaving")
		}

		_, err = b.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
			ChatID:          cc.Message.Chat.ID,
			InlineMessageID: cc.InlineMessageID,
			MessageID:       cc.Message.ID,
			ReplyMarkup: tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(btnLanguage),
				tg.NewInlineKeyboardRow(btnAutosaving),
			),
		})
	case "toggle_autosaving":
		user.AutoSaving = !user.AutoSaving
		if err = e.store.UpdateUser(user); err != nil {
			answer.Text = "🤷🏻‍♂️ " + p.Sprintf("Unknown error")
			_, err = b.AnswerCallbackQuery(answer)
			return err
		}
		fallthrough
	case "goto_autosaving":
		btnState := tg.NewInlineKeyboardButton(
			"☑️ "+strings.ToTitle(p.Sprintf("disabled")),
			"toggle_autosaving",
		)
		if user.AutoSaving {
			btnState.Text = "✅ " + strings.ToTitle(p.Sprintf("enabled"))
		}

		_, err = b.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
			ChatID:          cc.Message.Chat.ID,
			InlineMessageID: cc.InlineMessageID,
			MessageID:       cc.Message.ID,
			ReplyMarkup: tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(btnState),
				tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButton("↩️ "+p.Sprintf("Go back"), "goto_settings"),
				),
			),
		})
	case "set_language":
		user.LanguageCode = data[1]
		if err = e.store.UpdateUser(user); err != nil {
			_, err = b.AnswerCallbackQuery(answer)
			return err
		}
		p = message.NewPrinter(language.Make(user.LanguageCode))
		fallthrough
	case "goto_languages":
		markup := tg.NewInlineKeyboardMarkup()
		for _, l := range message.DefaultCatalog.Languages() {
			base, _ := l.Base()
			if user.LanguageCode == base.String() {
				continue
			}

			markup.InlineKeyboard = append(
				markup.InlineKeyboard, tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButton(
					message.NewPrinter(l).Sprintf("language_flag")+" "+strings.ToTitle(display.Self.Name(l)),
					"set_language:"+base.String(),
				)),
			)
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton("↩️ "+p.Sprintf("Go back"), "goto_settings"),
		))

		_, err = b.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParameters{
			ChatID:          cc.Message.Chat.ID,
			InlineMessageID: cc.InlineMessageID,
			MessageID:       cc.Message.ID,
			ReplyMarkup:     markup,
		})
	default:
		answer.Text = "🤷🏻‍♂️ " + p.Sprintf("Unknown error")
	}
	_, err = b.AnswerCallbackQuery(answer)
	return err
}
