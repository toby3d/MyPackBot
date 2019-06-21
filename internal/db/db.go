package db

import (
	"os"

	bolt "github.com/etcd-io/bbolt"
)

func New(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}

	if err = AutoMigrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) (err error) {
		if _, err = tx.CreateBucketIfNotExists([]byte("users")); err != nil {
			return err
		}

		if _, err = tx.CreateBucketIfNotExists([]byte("stickers")); err != nil {
			return err
		}

		if _, err = tx.CreateBucketIfNotExists([]byte("users_stickers")); err != nil {
			return err
		}

		return nil
	})
}
