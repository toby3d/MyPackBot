package store

import (
	"strconv"
	"strings"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/json-iterator/go"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	store "gitlab.com/toby3d/mypackbot/internal/model/store"
)

type Store struct {
	conn     *bolt.DB
	users    store.UsersManager
	stickers store.StickersManager
}

func NewStore(conn *bolt.DB) *Store {
	return &Store{
		conn:     conn,
		users:    NewUsersStore(conn),
		stickers: NewStickersStore(conn),
	}
}

func (store *Store) Users() store.UsersManager { return store.users }

func (store *Store) Stickers() store.StickersManager { return store.stickers }

func (store *Store) AddSticker(u *model.User, s *model.Sticker) (err error) {
	var us *model.UserSticker

	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}

	if us != nil {
		return common.ErrUserStickerExist
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

	err = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsersStickers).ForEach(func(key, val []byte) error {
			item := new(model.UserSticker)

			if err := json.ConfigFastest.Unmarshal(val, item); err != nil {
				return err
			}

			if item.UserID != u.ID || item.StickerID != s.ID {
				return nil
			}

			us = item

			return nil
		})
	})

	return us, err
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

	return stickers, count
}

func (store *Store) RemoveSticker(u *model.User, s *model.Sticker) (err error) {
	var us *model.UserSticker

	if us, err = store.GetSticker(u, s); err != nil {
		return err
	}

	if us == nil {
		return common.ErrUserStickerNotExist
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketUsersStickers)

		return bkt.ForEach(func(key, val []byte) error {
			item := new(model.UserSticker)

			if err := json.ConfigFastest.Unmarshal(val, &item); err != nil {
				return err
			}

			if item.UserID != us.UserID || item.StickerID != us.StickerID {
				return nil
			}

			return bkt.Delete(key)
		})
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
