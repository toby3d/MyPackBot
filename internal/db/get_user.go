package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
	"golang.org/x/text/language"
)

func (db *DB) GetUser(u *User) (*User, error) {
	if db == nil || db.db == nil {
		return nil, ErrDatabaseClosed
	}
	if u == nil {
		return nil, ErrNotFound
	}

	err := db.db.Batch(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(bucketUsers)
		if err != nil {
			return err
		}

		_, err = getUser(users, u)
		return err
	})
	return u, err
}

func getUser(b *bolt.Bucket, u *User) (*bolt.Bucket, error) {
	user := b.Bucket(strconv.AppendInt(nil, int64(u.ID), 10))
	if user == nil {
		return createUser(b, u)
	}

	u.Language = language.Make(string(user.Get(keyLanguage)))
	u.Sort = string(user.Get(keySort))
	u.State = string(user.Get(keyState))

	var err error
	if u.Hits, err = strconv.Atoi(string(user.Get(keyHits))); err != nil {
		return user, err
	}

	err = user.Tx().ForEach(func(name []byte, b *bolt.Bucket) error {
		set := &Set{
			ID:   string(name),
			User: u,
		}

		if _, err := getSet(user, set); err != nil {
			return err
		}

		u.Sets = append(u.Sets, set)
		return nil
	})

	return user, err
}
