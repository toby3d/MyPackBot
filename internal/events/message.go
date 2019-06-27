package events

import (
	"time"

	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Command struct {
	ID          string // start
	Description string // Just start the bot from scratch
	Message     string // Hello, %s! This bot...
}

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
	case m.IsCommand():
		switch {
		case m.IsCommandEqual(tg.CommandStart):
			text := p.Sprintf("Hello, %s!\n", m.From.FirstName)
			reply := tg.NewMessage(m.Chat.ID, text)
			_, err = b.SendMessage(reply)
		case m.IsCommandEqual(tg.CommandHelp):
			text := p.Sprintf("/start - just start\n/help - show this message\n/settings - get your current settings", m.From.FirstName)
			reply := tg.NewMessage(m.Chat.ID, text)
			_, err = b.SendMessage(reply)
		case m.IsCommandEqual(tg.CommandSettings):
			text := p.Sprintf("Here your settings")
			reply := tg.NewMessage(m.Chat.ID, text)

			btnLanguage := tg.NewInlineKeyboardButton(
				p.Sprintf("language_flag")+" "+p.Sprintf("Language"), "goto_languages",
			)

			btnAutosaving := tg.NewInlineKeyboardButton("☑️ "+p.Sprintf("Autosaving"), "goto_autosaving")
			if user.AutoSaving {
				btnAutosaving.Text = "✅ " + p.Sprintf("Autosaving")
			}

			reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(btnLanguage),
				tg.NewInlineKeyboardRow(btnAutosaving),
			)

			_, err = b.SendMessage(reply)
			// TODO: case m.IsCommandEqual("demo"):
		}
	case m.IsSticker(): // Input sticker
		// Get exists or create new sticker in database
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

		set := e.store.GetSet(sticker.SetName)
		if sticker.InSet() && set == nil {
			// This is a new set for database, get Title and stickers list from Telegram servers
			tgSet, err := b.GetStickerSet(m.Sticker.SetName)
			if err != nil {
				return err
			}

			// Rewrite current empty set on created new
			if set, err = e.store.GetOrCreateSet(&models.Set{
				Name:  tgSet.Name,
				Title: tgSet.Title,
			}); err != nil {
				return err
			}

			go func() { // Import all avaliable stickers in background
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

		// Create keyboard markup for current sticker
		markup := tg.NewInlineKeyboardMarkup()

		// This sticker is already exists in user pack
		if e.store.GetUserSticker(user.ID, sticker.ID) != nil {
			markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButton("🔥 "+p.Sprintf("Remove this sticker"), "remove_sticker"),
			))
			if sticker.InSet() && set != nil {
				markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButton(
						"🔥 "+p.Sprintf("Remove %s set", set.Title), "remove_set",
					),
				))
			}
		} else {
			// Automaticly add sticker into user pack if autosaving enabled
			if user.AutoSaving {
				if err = e.store.AddSticker(user.ID, sticker.ID); err != nil {
					return err
				}

				markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButton(
						"🔥 "+p.Sprintf("Remove this sticker"), "remove_sticker",
					),
				))
			} else {
				markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButton("❤️ "+p.Sprintf("Add this sticker"), "add_sticker"),
				))
			}
			if sticker.InSet() {
				markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButton(
						"💕 "+p.Sprintf("Add %s set", set.Title), "add_set",
					),
				))
			}
		}

		// Delete original sticker message for use keyboard workaround
		if _, err = b.DeleteMessage(m.Chat.ID, m.ID); err != nil {
			return err
		}

		// Send same sticker with keyboard
		_, err = b.SendSticker(&tg.SendStickerParameters{
			ChatID:      m.Chat.ID,
			Sticker:     sticker.ID,
			ReplyMarkup: markup,
		})
	}
	return err
}
