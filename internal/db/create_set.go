package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) CreateSet(s *Set) error {
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

		if _, err = createSet(user, s); err != nil {
			return err
		}
		return nil
	})
}

func createSet(b *bolt.Bucket, s *Set) (*bolt.Bucket, error) {
	set, err := b.CreateBucket([]byte(s.ID))
	if err == ErrAlreadyExist {
		return getSet(b, s)
	}
	if err != nil {
		return nil, err
	}

	if err = set.Put(keyHits, strconv.AppendInt(nil, int64(s.Hits), 10)); err != nil {
		return set, err
	}
	if err = set.Put(keyIsFavorite, strconv.AppendBool(nil, s.IsFavorite)); err != nil {
		return set, err
	}
	if err = set.Put(keySort, valSortHits); err != nil {
		return set, err
	}

	for i := range s.Stickers {
		if _, err = createSticker(set, s.Stickers[i]); err != nil {
			return nil, err
		}
	}

	return set, nil
}
