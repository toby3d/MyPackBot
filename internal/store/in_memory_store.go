package store

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"gitlab.com/toby3d/mypackbot/internal/model"
	store "gitlab.com/toby3d/mypackbot/internal/model/store"
)

type InMemoryStore struct {
	users    store.UsersManager
	stickers store.StickersManager

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

func (store *InMemoryStore) Users() store.UsersManager { return store.users }

func (store *InMemoryStore) Stickers() store.StickersManager { return store.stickers }

func (store *InMemoryStore) AddSticker(u *model.User, s *model.Sticker, emoji string) (err error) {
	var us *model.UserSticker
	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}
	if us != nil {
		return errors.New("sticker already added to this user")
	}

	store.mutex.Lock()
	store.userStickers = append(store.userStickers, &model.UserSticker{
		UserID:    u.ID,
		StickerID: s.ID,
		Emoji:     emoji,
		Hits:      0,
	})
	sort.Slice(store.userStickers, func(i, j int) bool {
		return store.userStickers[i].UserID < store.userStickers[j].UserID ||
			store.userStickers[i].StickerID < store.userStickers[j].StickerID
	})
	store.mutex.Unlock()

	return err
}

func (store *InMemoryStore) GetSticker(u *model.User, s *model.Sticker) (us *model.UserSticker, err error) {
	if u, err = store.Users().GetOrCreate(u); err != nil {
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
	var count int
	stickers := make(model.Stickers, 0, limit)

	store.mutex.RLock()
	for i := range store.userStickers {
		if store.userStickers[i].UserID != u.ID ||
			(query != "" && !strings.ContainsAny(store.userStickers[i].Emoji, query)) {
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

func (store *InMemoryStore) GetStickersSet(u *model.User, offset, limit int, setName string) (model.Stickers, int) {
	var count int
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
		return errors.New("sticker already removed in this user")
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
		si := store.userStickers[i]
		sj := store.userStickers[j]
		return store.Stickers().Get(si.StickerID).SetName < store.Stickers().Get(sj.StickerID).SetName ||
			si.UserID < sj.UserID || si.StickerID < sj.StickerID
	})
	store.mutex.Unlock()

	return err
}

func (store *InMemoryStore) HitSticker(u *model.User, s *model.Sticker) (err error) {
	var us *model.UserSticker
	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}
	if us == nil {
		return errors.New("sticker not exists in this user")
	}

	return store.stickers.Hit(us.StickerID)
}
