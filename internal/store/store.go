package store

import (
	"sort"
	"strconv"
	"strings"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/json-iterator/go"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	"golang.org/x/xerrors"
)

type Store struct {
	conn     *bolt.DB
	users    users.Manager
	stickers stickers.Manager
}

func NewStore(conn *bolt.DB) *Store {
	return &Store{
		conn:     conn,
		users:    NewUsersStore(conn),
		stickers: NewStickersStore(conn),
	}
}

func (store *Store) AddSticker(u *model.User, s *model.Sticker) (err error) {
	var us *model.UserSticker

	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}

	if us != nil {
		return model.ErrUserStickerExist
	}

	src, err := json.ConfigFastest.Marshal(&model.UserSticker{
		CreatedAt: time.Now().UTC().Unix(),
		Emojis:    s.Emoji,
		SetName:   s.SetName,
		StickerID: s.ID,
		UserID:    u.ID,
	})
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketUsersStickers)

		id, err := bkt.NextSequence()
		if err != nil {
			return err
		}

		return bkt.Put([]byte(strconv.FormatUint(id, 10)), src)
	})
}

func (store *Store) AddStickersSet(u *model.User, setName string) (err error) {
	if u, err = store.users.GetOrCreate(u); err != nil {
		return err
	}

	set, _ := store.stickers.GetSet(setName)
	for _, s := range set {
		_ = store.AddSticker(u, s)
	}

	return err
}

func (store *Store) GetSticker(u *model.User, s *model.Sticker) (*model.UserSticker, error) {
	var err error

	if u, err = store.users.GetOrCreate(u); err != nil {
		return nil, err
	}

	if s, err = store.stickers.GetOrCreate(s); err != nil {
		return nil, err
	}

	var us *model.UserSticker

	if err = store.conn.View(func(tx *bolt.Tx) (err error) {
		if err = tx.Bucket(common.BucketUsersStickers).ForEach(func(key, val []byte) (err error) {
			item := new(model.UserSticker)

			if err = json.ConfigFastest.Unmarshal(val, item); err != nil {
				return err
			}

			if item.UserID != u.ID || item.StickerID != s.ID {
				return nil
			}

			us = item

			return model.ErrForEachStop
		}); err != nil && !xerrors.Is(err, model.ErrForEachStop) {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return us, nil
}

func (store *Store) GetStickersList(u *model.User, offset, limit int, query string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0, limit)
	_ = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsersStickers).ForEach(func(key, val []byte) (err error) {
			us := new(model.UserSticker)

			if err = json.ConfigFastest.Unmarshal(val, us); err != nil {
				return err
			}

			if us.UserID != u.ID || (query != "" && !strings.ContainsAny(us.Emojis, query)) {
				return nil
			}

			count++

			if (offset != 0 && count <= offset) || count > offset+limit {
				return nil
			}

			s := new(model.Sticker)
			src := tx.Bucket(common.BucketStickers).Get([]byte(us.StickerID))

			if err = json.ConfigFastest.Unmarshal(src, s); err != nil {
				return err
			}

			stickers = append(stickers, s)

			return nil
		})
	})

	sort.Slice(stickers, func(i, j int) bool {
		return stickers[i].SetName < stickers[j].SetName || stickers[i].UpdatedAt < stickers[j].UpdatedAt
	})

	return stickers, count
}

func (store *Store) GetStickersSet(u *model.User, offset, limit int, setName string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0, limit)

	_ = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsersStickers).ForEach(func(key, val []byte) (err error) {
			us := new(model.UserSticker)

			if err = json.ConfigFastest.Unmarshal(val, us); err != nil {
				return err
			}

			if us.UserID != u.ID {
				return nil
			}

			s := new(model.Sticker)
			src := tx.Bucket(common.BucketStickers).Get([]byte(us.StickerID))

			if err = json.ConfigFastest.Unmarshal(src, s); err != nil {
				return err
			}

			if !strings.EqualFold(s.SetName, setName) {
				return nil
			}

			count++

			if count <= offset || count > limit {
				return nil
			}

			stickers = append(stickers, s)

			return nil
		})
	})

	sort.Slice(stickers, func(i, j int) bool {
		return stickers[i].UpdatedAt < stickers[j].UpdatedAt
	})

	return stickers, count
}

func (store *Store) RemoveSticker(u *model.User, s *model.Sticker) (err error) {
	var us *model.UserSticker

	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}

	if us == nil {
		return model.ErrUserStickerNotExist
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsersStickers)

		if err = bkt.ForEach(func(key, val []byte) (err error) {
			item := new(model.UserSticker)

			if err = json.ConfigFastest.Unmarshal(val, &item); err != nil {
				return err
			}

			if item.UserID != us.UserID || item.StickerID != us.StickerID {
				return nil
			}

			if err = bkt.Delete(key); err != nil {
				return err
			}

			return model.ErrForEachStop
		}); err != nil && !xerrors.Is(err, model.ErrForEachStop) {
			return err
		}

		return nil
	})
}

func (store *Store) RemoveStickersSet(u *model.User, setName string) (err error) {
	if u, err = store.users.GetOrCreate(u); err != nil {
		return err
	}

	set, _ := store.stickers.GetSet(setName)
	for _, s := range set {
		_ = store.RemoveSticker(u, s)
	}

	return err
}
