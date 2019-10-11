package event

import (
	"gitlab.com/toby3d/mypackbot/internal/middleware"
	"gitlab.com/toby3d/mypackbot/internal/model"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

func (event *Event) Message(m *tg.Message) error {
	switch {
	case m.IsCommand():
		return event.Commands(m)
	case m.IsSticker():
		return event.Stickers(m)
	}
	return nil
}

func (event *Event) Commands(m *tg.Message) error {
	u, err := middleware.GetUser(event.store, m.From, m.Date)
	if err != nil {
		return err
	}

	p := middleware.GetPrinter(u.LanguageCode)
	switch {
	case m.IsCommandEqual(tg.CommandStart):
		return event.StartCommand(p, m)
	case m.IsCommandEqual(tg.CommandHelp):
		return event.HelpCommand(p, m)
	case m.IsCommandEqual(tg.CommandSettings):
		return event.SettingsCommand(p, m)
	default:
		return event.UnknownCommand(p, m)
	}
}

func (event *Event) StartCommand(p *message.Printer, m *tg.Message) (err error) {
	reply := tg.NewMessage(m.Chat.ID, p.Sprintf("start__text", m.From.FullName()))
	reply.ReplyToMessageID = m.ID
	_, err = event.bot.SendMessage(reply)
	return err
}

func (event *Event) HelpCommand(p *message.Printer, m *tg.Message) (err error) {
	reply := tg.NewMessage(m.Chat.ID, p.Sprintf("help__text"))
	reply.ReplyToMessageID = m.ID
	_, err = event.bot.SendMessage(reply)
	return err
}

func (event *Event) SettingsCommand(p *message.Printer, m *tg.Message) (err error) {
	reply := tg.NewMessage(m.Chat.ID, p.Sprintf("settings-command__text"))
	reply.ReplyToMessageID = m.ID
	_, err = event.bot.SendMessage(reply)
	return err
}

func (event *Event) UnknownCommand(p *message.Printer, m *tg.Message) (err error) {
	reply := tg.NewMessage(m.Chat.ID, p.Sprintf("unknown-command__text"))
	reply.ReplyToMessageID = m.ID
	_, err = event.bot.SendMessage(reply)
	return err
}

func (event *Event) Stickers(m *tg.Message) error {
	u, err := middleware.GetUser(event.store, m.From, m.Date)
	if err != nil {
		return err
	}

	s, err := event.store.Stickers().GetOrCreate(&model.Sticker{
		ID:         m.Sticker.FileID,
		IsAnimated: m.Sticker.IsAnimated,
		SetName:    m.Sticker.SetName,
		CreatedAt:  m.Date,
	})
	if err != nil {
		return err
	}

	us, err := event.store.GetSticker(u, s)
	if err != nil {
		return err
	}

	p := middleware.GetPrinter(u.LanguageCode)
	markup := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-single"), "add:single"),
	))
	if s.SetName != "" {
		markup.InlineKeyboard[0] = append(
			markup.InlineKeyboard[0],
			tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_add-set"), "add:set"),
		)
	}
	if us != nil {
		markup = tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-single"), "remove:single"),
		))
		if s.SetName != "" {
			markup.InlineKeyboard[0] = append(
				markup.InlineKeyboard[0],
				tg.NewInlineKeyboardButton(p.Sprintf("sticker__button_remove-set"), "remove:set"),
			)
		}
	}

	reply := tg.NewMessage(m.Chat.ID, p.Sprintf("sticker__text"))
	reply.ReplyToMessageID = m.ID
	reply.ReplyMarkup = markup
	_, err = event.bot.SendMessage(reply)
	return err
}
