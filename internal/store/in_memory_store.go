package store

import (
	"sort"
	"strings"
	"sync"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
)

type InMemoryStore struct {
	users    users.Manager
	stickers stickers.Manager

	mutex        sync.RWMutex
	userStickers model.UserStickers
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users:    NewInMemoryUsersStore(),
		stickers: NewInMemoryStickersStore(),
		mutex:    sync.RWMutex{},
	}
}

func (store *InMemoryStore) Users() users.Manager { return store.users }

func (store *InMemoryStore) Stickers() stickers.Manager { return store.stickers }

func (store *InMemoryStore) AddSticker(u *model.User, s *model.Sticker) (err error) {
	var us *model.UserSticker

	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}

	if us != nil {
		return model.ErrUserStickerExist
	}

	store.mutex.Lock()
	store.userStickers = append(store.userStickers, &model.UserSticker{
		CreatedAt: time.Now().UTC().Unix(),
		Emojis:    s.Emoji,
		SetName:   s.SetName,
		StickerID: s.ID,
		UserID:    u.ID,
	})

	sort.Slice(store.userStickers, func(i, j int) bool {
		return store.userStickers[i].SetName < store.userStickers[j].SetName ||
			store.userStickers[i].CreatedAt < store.userStickers[j].CreatedAt
	})
	store.mutex.Unlock()

	return err
}

func (store *InMemoryStore) AddStickersSet(u *model.User, setName string) (err error) {
	if u, err = store.users.GetOrCreate(u); err != nil {
		return err
	}

	set, _ := store.stickers.GetSet(setName)
	for _, s := range set {
		_ = store.AddSticker(u, s)
	}

	return err
}

func (store *InMemoryStore) GetSticker(u *model.User, s *model.Sticker) (us *model.UserSticker, err error) {
	if u, err = store.users.GetOrCreate(u); err != nil {
		return nil, err
	}

	if s, err = store.Stickers().GetOrCreate(s); err != nil {
		return nil, err
	}

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	return store.userStickers.GetByID(u.ID, s.ID), err
}

func (store *InMemoryStore) GetStickersList(u *model.User, offset, limit int, query string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0, limit)

	store.mutex.RLock()
	for i := range store.userStickers {
		if store.userStickers[i].UserID != u.ID ||
			(query != "" && !strings.ContainsAny(store.userStickers[i].Emojis, query)) {
			continue
		}

		count++

		if (offset != 0 && count <= offset) || count > offset+limit {
			continue
		}

		stickers = append(stickers, store.Stickers().Get(store.userStickers[i].StickerID))
	}
	store.mutex.RUnlock()

	return stickers, count
}

func (store *InMemoryStore) GetStickersSet(u *model.User, offset, limit int, setName string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0, limit)

	store.mutex.RLock()
	for i := range store.userStickers {
		if store.userStickers[i].UserID != u.ID {
			continue
		}

		s := store.stickers.Get(store.userStickers[i].StickerID)

		if !strings.EqualFold(s.SetName, setName) {
			continue
		}

		count++

		if count < offset || count > limit {
			continue
		}

		stickers = append(stickers, store.Stickers().Get(store.userStickers[i].StickerID))
	}
	store.mutex.RUnlock()

	return stickers, count
}

func (store *InMemoryStore) RemoveSticker(u *model.User, s *model.Sticker) (err error) {
	var us *model.UserSticker

	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}

	if us == nil {
		return model.ErrUserStickerNotExist
	}

	store.mutex.Lock()
	for i := range store.userStickers {
		if store.userStickers[i].UserID != u.ID || store.userStickers[i].StickerID != s.ID {
			continue
		}

		store.userStickers = store.userStickers[:i+copy(store.userStickers[i:], store.userStickers[i+1:])]

		break
	}

	sort.Slice(store.userStickers, func(i, j int) bool {
		return store.userStickers[i].SetName < store.userStickers[j].SetName ||
			store.userStickers[i].CreatedAt < store.userStickers[j].CreatedAt
	})
	store.mutex.Unlock()

	return err
}

func (store *InMemoryStore) RemoveStickersSet(u *model.User, setName string) (err error) {
	if u, err = store.users.GetOrCreate(u); err != nil {
		return err
	}

	set, _ := store.stickers.GetSet(setName)
	for _, s := range set {
		_ = store.RemoveSticker(u, s)
	}

	return err
}
