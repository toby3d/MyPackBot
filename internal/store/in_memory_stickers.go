package store

import (
	"sort"
	"strings"
	"sync"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
)

type InMemoryStickersStore struct {
	mutex    sync.RWMutex
	stickers model.Stickers
}

func NewInMemoryStickersStore() *InMemoryStickersStore {
	return &InMemoryStickersStore{
		mutex:    sync.RWMutex{},
		stickers: make([]*model.Sticker, 0),
	}
}

func (store *InMemoryStickersStore) Create(s *model.Sticker) error {
	if store.Get(s.ID) != nil {
		return model.ErrStickerExist
	}

	if s.CreatedAt == 0 {
		s.CreatedAt = time.Now().UTC().Unix()
	}

	store.mutex.Lock()
	store.stickers = append(store.stickers, s)

	sort.Slice(store.stickers, func(i, j int) bool {
		return store.stickers[i].CreatedAt < store.stickers[j].CreatedAt
	})
	store.mutex.Unlock()

	return nil
}

func (store *InMemoryStickersStore) Get(sid string) *model.Sticker {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	return store.stickers.GetByID(sid)
}

func (store *InMemoryStickersStore) GetList(offset, limit int, query string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0, limit)

	store.mutex.RLock()
	for i := range store.stickers {
		if query != "" && !strings.ContainsAny(store.stickers[i].Emoji, query) {
			continue
		}

		count++

		if count <= offset || count > offset+limit {
			continue
		}

		stickers = append(stickers, store.stickers[i])
	}
	store.mutex.RUnlock()

	return stickers, count
}

func (store *InMemoryStickersStore) GetSet(name string) (model.Stickers, int) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	return store.stickers.GetSet(name)
}

func (store *InMemoryStickersStore) Update(s *model.Sticker) error {
	if store.Get(s.ID) == nil {
		return store.Create(s)
	}

	store.mutex.Lock()
	for i := range store.stickers {
		if store.stickers[i].ID != s.ID {
			continue
		}

		store.stickers[i] = s
	}
	store.mutex.Unlock()

	return nil
}

func (store *InMemoryStickersStore) Remove(sid string) error {
	if store.Get(sid) == nil {
		return model.ErrStickerNotExist
	}

	store.mutex.Lock()
	for i := range store.stickers {
		if store.stickers[i].ID != sid {
			continue
		}

		store.stickers = store.stickers[:i+copy(store.stickers[i:], store.stickers[i+1:])]

		break
	}

	sort.Slice(store.stickers, func(i, j int) bool {
		return store.stickers[i].CreatedAt < store.stickers[j].CreatedAt
	})
	store.mutex.Unlock()

	return nil
}

func (store *InMemoryStickersStore) GetOrCreate(s *model.Sticker) (*model.Sticker, error) {
	if sticker := store.Get(s.ID); sticker != nil {
		return sticker, nil
	}

	if err := store.Create(s); err != nil {
		return nil, err
	}

	return store.Get(s.ID), nil
}
