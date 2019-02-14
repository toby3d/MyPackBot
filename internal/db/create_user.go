package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
)

func (db *DB) CreateUser(u *User) error {
	if db == nil || db.db == nil {
		return ErrDatabaseClosed
	}
	if u == nil {
		return ErrNotFound
	}

	return db.db.Batch(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(bucketUsers)
		if err != nil {
			return err
		}

		_, err = createUser(users, u)
		return err
	})
}

func createUser(b *bolt.Bucket, u *User) (*bolt.Bucket, error) {
	user, err := b.CreateBucket(strconv.AppendInt(nil, int64(u.ID), 10))
	if err == ErrAlreadyExist {
		return getUser(b, u)
	}
	if err != nil {
		return nil, err
	}

	lang, _ := u.Language.Base()
	if err = user.Put(keyLanguage, []byte(lang.String())); err != nil {
		return user, err
	}
	if err = user.Put(keyHits, strconv.AppendInt(nil, int64(u.Hits), 10)); err != nil {
		return user, err
	}
	if err = user.Put(keySort, valSortHits); err != nil {
		return user, err
	}

	for i := range u.Sets {
		if _, err := createSet(user, u.Sets[i]); err != nil {
			return user, err
		}
	}

	return user, nil
}
