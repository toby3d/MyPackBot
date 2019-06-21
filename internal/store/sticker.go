package store

import (
	"strings"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/pquerna/ffjson/ffjson"
	"gitlab.com/toby3d/mypackbot/internal/models"
)

type StickerStore struct {
	db *bolt.DB
}

func NewStickerStore(db *bolt.DB) *StickerStore {
	return &StickerStore{db: db}
}

func (ss *StickerStore) GetByID(sid string) (*models.Sticker, error) {
	var s models.Sticker
	err := ss.db.View(func(tx *bolt.Tx) error {
		src := tx.Bucket([]byte("stickers")).Get([]byte(sid))
		if src == nil {
			return nil
		}
		return json.UnmarshalFast(src, &s)
	})
	return &s, err
}

func (ss *StickerStore) GetByUserID(uid int, offset, limit int) ([]models.Sticker, int, error) {
	var stickers []models.Sticker
	var count int
	err := ss.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("users_stickers")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var us models.UsersStickers
			if err := json.UnmarshalFast(v, &us); err != nil {
				return err
			}

			if us.UserID != uid {
				continue
			}
			sticker, err := ss.GetByID(us.StickerID)
			if err != nil {
				return err
			}
			if sticker == nil {
				continue
			}
			count++

			if offset != -1 && offset > 0 {
				offset--
				continue
			}

			if limit != -1 {
				if limit > 0 {
					limit--
				} else {
					continue
				}
			}
			stickers = append(stickers, *sticker)
		}
		return nil
	})
	return stickers, count, err
}

func (ss *StickerStore) Create(s *models.Sticker) error {
	src, err := json.MarshalFast(s)
	if err != nil {
		return err
	}

	return ss.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("stickers")).Put([]byte(s.ID), src)
	})
}

func (ss *StickerStore) Update(s *models.Sticker) error {
	src, err := json.MarshalFast(s)
	if err != nil {
		return err
	}

	return ss.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("stickers")).Put([]byte(s.ID), src)
	})
}

func (ss *StickerStore) Delete(s *models.Sticker) error {
	return ss.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("stickers")).Delete([]byte(s.ID))
	})
}

func (ss *StickerStore) List(offset, limit int) ([]models.Sticker, int, error) {
	var stickers []models.Sticker
	var count int
	err := ss.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("stickers")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			count++
			if offset != -1 && offset > 0 {
				offset--
				continue
			}
			if limit != -1 {
				if limit > 0 {
					limit--
				} else {
					continue
				}
			}

			var s models.Sticker
			if err := json.UnmarshalFast(v, &s); err != nil {
				return err
			}
			stickers = append(stickers, s)
		}
		return nil
	})
	return stickers, count, err
}

func (ss *StickerStore) ListByEmoji(emoji string, offset, limit int) ([]models.Sticker, int, error) {
	var stickers []models.Sticker
	var count int
	err := ss.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("stickers")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var s models.Sticker
			if err := json.UnmarshalFast(v, &s); err != nil {
				return err
			}

			if !strings.ContainsAny(s.Emoji, emoji) {
				continue
			}
			count++

			if offset != -1 && offset > 0 {
				offset--
				continue
			}

			if limit != -1 {
				if limit > 0 {
					limit--
				} else {
					continue
				}
			}
			stickers = append(stickers, s)
		}
		return nil
	})
	return stickers, count, err
}

func (ss *StickerStore) GetSet(setName string, offset, limit int) ([]models.Sticker, int, error) {
	var stickers []models.Sticker
	var count int
	err := ss.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("stickers")).Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var s models.Sticker
			if err := json.UnmarshalFast(v, &s); err != nil {
				return err
			}

			if !strings.EqualFold(s.SetName, setName) {
				continue
			}
			count++

			if offset != -1 && offset > 0 {
				offset--
				continue
			}

			if limit != -1 {
				if limit > 0 {
					limit--
				} else {
					continue
				}
			}
			stickers = append(stickers, s)
		}
		return nil
	})
	return stickers, count, err
}
