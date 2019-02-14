package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) GetSticker(s *Sticker) (*Sticker, error) {
	if db == nil || db.db == nil {
		return nil, ErrDatabaseClosed
	}
	if s == nil {
		return nil, ErrNotFound
	}

	err := db.db.Batch(func(tx *bolt.Tx) error {
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

		_, err = getSticker(set, s)
		return err
	})
	return s, err
}

func getSticker(b *bolt.Bucket, s *Sticker) (*bolt.Bucket, error) {
	sticker := b.Bucket([]byte(s.ID))
	if sticker == nil {
		return createSticker(b, s)
	}

	s.Emoji = string(sticker.Get(keyEmoji))

	var err error
	if s.Hits, err = strconv.Atoi(string(sticker.Get(keyHits))); err != nil {
		return sticker, err
	}

	if s.IsFavorite, err = strconv.ParseBool(string(sticker.Get(keyIsFavorite))); err != nil {
		return sticker, err
	}

	return sticker, nil
}
