package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) CreateSticker(s *Sticker) error {
	if db == nil || db.db == nil {
		return ErrDatabaseClosed
	}
	if s == nil {
		return ErrNotFound
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

		_, err = createSticker(set, s)
		return err
	})
}

func createSticker(b *bolt.Bucket, s *Sticker) (*bolt.Bucket, error) {
	sticker, err := b.CreateBucket([]byte(s.ID))
	if err == ErrAlreadyExist {
		return getSticker(b, s)
	}
	if err != nil {
		return nil, err
	}

	if err = sticker.Put(keyEmoji, []byte(s.Emoji)); err != nil {
		return sticker, err
	}
	if err = sticker.Put(keyHits, strconv.AppendInt(nil, int64(s.Hits), 10)); err != nil {
		return sticker, err
	}
	if err = sticker.Put(keyIsFavorite, strconv.AppendBool(nil, s.IsFavorite)); err != nil {
		return sticker, err
	}

	return sticker, err
}
