package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) UpdateSticker(s *Sticker) error {
	if db == nil || db.db == nil {
		return ErrDatabaseClosed
	}

	return db.db.Batch(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(bucketUsers)
		if err != nil {
			return err
		}

		user, err := getUser(users, s.User)
		if err != nil {
			return err
		}

		set, err := getSet(user, s.Set)
		if err != nil {
			return err
		}

		sticker, err := getSticker(set, s)
		if err != nil {
			return err
		}

		if s.Emoji != "" {
			if err = sticker.Put(keyEmoji, []byte(s.Emoji)); err != nil {
				return err
			}
		}

		if s.Hits > 0 {
			if err = sticker.Put(keyHits, strconv.AppendInt(nil, int64(s.Hits), 10)); err != nil {
				return err
			}
		}

		if err = sticker.Put(
			keyIsFavorite,
			strconv.AppendBool(nil, s.IsFavorite),
		); err != nil {
			return err
		}

		return nil
	})
}
