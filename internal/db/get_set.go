package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) GetSet(s *Set) (*Set, error) {
	if db == nil || db.db == nil {
		return nil, ErrDatabaseClosed
	}
	if s == nil {
		return nil, ErrNotFound
	}

	err := db.db.Update(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(bucketUsers)
		if err != nil {
			return err
		}

		user, err := getUser(users, s.User)
		if err != nil {
			return err
		}

		_, err = getSet(user, s)
		return err
	})
	return s, err
}

func getSet(b *bolt.Bucket, s *Set) (*bolt.Bucket, error) {
	set := b.Bucket([]byte(s.ID))
	if set == nil {
		return createSet(b, s)
	}

	var err error
	if s.Hits, err = strconv.Atoi(string(set.Get(keyHits))); err != nil {
		return set, err
	}

	if s.IsFavorite, err = strconv.ParseBool(string(set.Get(keyIsFavorite))); err != nil {
		return set, err
	}

	err = set.Tx().ForEach(func(name []byte, b *bolt.Bucket) error {
		var sticker Sticker
		sticker.ID = string(name)
		sticker.Set = s

		if _, err := getSticker(set, &sticker); err != nil {
			return err
		}

		s.Stickers = append(s.Stickers, &sticker)
		return nil
	})

	return set, err
}
