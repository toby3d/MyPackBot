package events

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (e *Events) Message(b *tg.Bot, m *tg.Message) error {
	if !m.Chat.IsPrivate() || !b.IsMessageToMe(m) {
		return nil
	}

	user, err := e.store.GetOrCreateUser(&models.User{
		ID:           m.From.ID,
		LanguageCode: m.From.LanguageCode,
		StartedAt:    m.Time().UTC().Unix(),
		AutoSaving:   true,
	})
	if err != nil {
		return err
	}
	p := message.NewPrinter(language.Make(user.LanguageCode))

	switch {
	case m.IsSticker():
		sticker, err := e.store.GetOrCreateSticker(&models.Sticker{
			Model: models.Model{
				ID:      m.Sticker.FileID,
				SavedAt: m.Time().UTC().Unix(),
			},
			Emoji:   m.Sticker.Emoji,
			SetName: m.Sticker.SetName,
		})
		if err != nil {
			return err
		}
		if err = e.store.AddSticker(user.ID, sticker.ID); err != nil {
			return err
		}

		reply := tg.NewMessage(m.Chat.ID, p.Sprintf("👍🏻 This sticker has been saved!"))
		reply.ReplyToMessageID = m.ID
		reply.ParseMode = tg.StyleMarkdown
		reply.DisableWebPagePreview = true

		if m.Sticker.InSet() {
			set := e.store.GetSet(m.Sticker.SetName)
			if set == nil {
				tgSet, err := b.GetStickerSet(m.Sticker.SetName)
				if err != nil {
					return err
				}

				if set, err = e.store.GetOrCreateSet(&models.Set{
					Name:  tgSet.Name,
					Title: tgSet.Title,
				}); err != nil {
					return err
				}

				go func() {
					for _, tgSticker := range tgSet.Stickers {
						tgSticker := tgSticker
						if _, err := e.store.GetOrCreateSticker(&models.Sticker{
							Model: models.Model{
								ID:      tgSticker.FileID,
								SavedAt: time.Now().UTC().Unix(),
							},
							Emoji:   tgSticker.Emoji,
							SetName: set.Name,
						}); err != nil {
							continue
						}
					}
				}()
			}

			reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(tg.NewInlineKeyboardButton(
					p.Sprintf("📥 Import all %s stickers", set.Title), "set:import:"+set.Name,
				)),
			)
			reply.Text = p.Sprintf(
				"This sticker from [%s](https://t.me/addstickers/%s) set has beed added!",
				set.Title, set.Name,
			)
		}

		_, err = b.SendMessage(reply)
		return err
	}
	return nil
}
