package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) UpdateSet(s *Set) error {
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

		set, err := getSet(user, s)
		if err != nil {
			return err
		}

		if s.Hits > 0 {
			if err = set.Put(
				keyHits,
				strconv.AppendInt(nil, int64(s.Hits), 10),
			); err != nil {
				return err
			}
		}

		if err = set.Put(
			keyIsFavorite,
			strconv.AppendBool(nil, s.IsFavorite),
		); err != nil {
			return err
		}

		return nil
	})
}
