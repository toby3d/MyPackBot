package db

import (
	"os"

	bolt "github.com/etcd-io/bbolt"
	"gitlab.com/toby3d/mypackbot/internal/common"
)

func Open(path string) (*bolt.DB, error) {
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
		for _, bkt := range [][]byte{
			common.BucketUsers,
			common.BucketStickers,
			common.BucketUsersStickers,
		} {
			if _, err = tx.CreateBucketIfNotExists(bkt); err != nil {
				return err
			}
		}
		return nil
	})
}
