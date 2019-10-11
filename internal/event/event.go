package event

import (
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
)

type Event struct {
	bot   *tg.Bot
	store store.Manager
}

func NewEvent(b *tg.Bot, s store.Manager) *Event {
	return &Event{
		bot:   b,
		store: s,
	}
}
