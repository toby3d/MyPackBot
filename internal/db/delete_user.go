package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) DeleteUser(u *User) error {
	if db == nil || db.db == nil {
		return ErrDatabaseClosed
	}

	return db.db.Batch(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(bucketUsers)
		if err != nil {
			return err
		}

		return users.DeleteBucket(strconv.AppendInt(nil, int64(u.ID), 10))
	})
}
