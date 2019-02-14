package db

import (
	"strconv"

	bolt "github.com/etcd-io/bbolt"
	"golang.org/x/text/language"
)

func (db *DB) UpdateUser(u *User) error {
	if db == nil || db.db == nil {
		return ErrDatabaseClosed
	}

	return db.db.Batch(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(bucketUsers)
		if err != nil {
			return err
		}

		user, err := getUser(users, u)
		if err != nil {
			return err
		}

		if u.Hits > 0 {
			if err = user.Put(keyHits, strconv.AppendInt(nil, int64(u.Hits), 10)); err != nil {
				return err
			}
		}

		if u.Language != language.Und {
			lang, _ := u.Language.Base()
			if err = user.Put(keyLanguage, []byte(lang.String())); err != nil {
				return err
			}
		}

		if u.Sort != "" {
			if err = user.Put(keySort, []byte(u.Sort)); err != nil {
				return err
			}
		}

		if u.State != "" {
			if err = user.Put(keyState, []byte(u.State)); err != nil {
				return err
			}
		}

		return nil
	})
}
