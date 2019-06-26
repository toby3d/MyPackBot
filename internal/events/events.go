package events

import (
	"gitlab.com/toby3d/mypackbot/internal/store"
)

type Events struct {
	store *store.Store
}

func New(s *store.Store) *Events {
	return &Events{store: s}
}
