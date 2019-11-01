package store

import (
	"errors"
	"strings"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/json-iterator/go"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
)

type StickersStore struct{ conn *bolt.DB }

func NewStickersStore(conn *bolt.DB) *StickersStore { return &StickersStore{conn: conn} }

func (store *StickersStore) Create(s *model.Sticker) error {
	if store.Get(s.ID) != nil {
		return errors.New("stickers already exists")
	}

	now := time.Now().UTC().Unix()

	if s.CreatedAt <= 0 {
		s.CreatedAt = now
	}

	if s.UpdatedAt <= 0 {
		s.UpdatedAt = now
	}

	if s.SetName == "" {
		s.SetName = common.SetNameUploaded
	}

	src, err := json.ConfigFastest.Marshal(s)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketStickers).Put([]byte(s.ID), src)
	})
}

func (store *StickersStore) Get(sid string) *model.Sticker {
	s := new(model.Sticker)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		src := tx.Bucket(common.BucketStickers).Get([]byte(sid))

		return json.ConfigFastest.Unmarshal(src, s)
	}); err != nil {
		return nil
	}

	return s
}

func (store *StickersStore) GetList(offset, limit int, query string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0, limit)
	_ = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketStickers).ForEach(func(key, val []byte) error {
			s := new(model.Sticker)

			if err := json.ConfigFastest.Unmarshal(val, s); err != nil {
				return err
			}

			if query != "" && !strings.ContainsAny(s.Emoji, query) {
				return nil
			}

			count++

			if count <= offset || count > offset+limit {
				return nil
			}

			stickers = append(stickers, s)

			return nil
		})
	})

	return stickers, count
}

func (store *StickersStore) GetSet(name string) (model.Stickers, int) {
	count := 0
	stickers := make(model.Stickers, 0)
	_ = store.conn.View(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketStickers).ForEach(func(key, val []byte) error {
			s := new(model.Sticker)

			if err := json.ConfigFastest.Unmarshal(val, s); err != nil {
				return err
			}

			if !strings.EqualFold(s.SetName, name) {
				return nil
			}

			count++

			stickers = append(stickers, s)

			return nil
		})
	})

	return stickers, count
}

func (store *StickersStore) Update(s *model.Sticker) error {
	if store.Get(s.ID) == nil {
		return store.Create(s)
	}

	if s.UpdatedAt <= 0 {
		s.UpdatedAt = time.Now().UTC().Unix()
	}

	if s.SetName == "" {
		s.SetName = common.SetNameUploaded
	}

	src, err := json.ConfigFastest.Marshal(s)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketStickers).Put([]byte(s.ID), src)
	})
}

func (store *StickersStore) Remove(sid string) error {
	if store.Get(sid) == nil {
		return errors.New("sticker already removed or not exists")
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketStickers)

		if err := bkt.Delete([]byte(sid)); err != nil {
			return err
		}

		bkt = tx.Bucket(common.BucketUsersStickers)

		return bkt.ForEach(func(key, val []byte) error {
			us := new(model.UserSticker)

			if err := json.Unmarshal(val, us); err != nil {
				return err
			}

			if us.StickerID != sid {
				return nil
			}

			return bkt.Delete(key)
		})
	})
}

func (store *StickersStore) GetOrCreate(s *model.Sticker) (*model.Sticker, error) {
	if sticker := store.Get(s.ID); sticker != nil {
		return sticker, nil
	}

	if err := store.Create(s); err != nil {
		return nil, err
	}

	return store.Get(s.ID), nil
}
