package db

import bolt "github.com/etcd-io/bbolt"

func (db *DB) DeleteSet(s *Set) error {
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

		return user.DeleteBucket([]byte(s.ID))
	})
}
