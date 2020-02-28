package db

import (
	"os"

	"github.com/timshannon/bolthold"
	"gitlab.com/toby3d/mypackbot/internal/common"
	bolt "go.etcd.io/bbolt"
)

func Open(path string) (*bolthold.Store, error) {
	db, err := bolthold.Open(path, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}

	err = db.Bolt().Update(func(tx *bolt.Tx) error {
		for i := range common.Buckets {
			if _, err := tx.CreateBucketIfNotExists(common.Buckets[i]); err == nil {
				continue
			}

			return err
		}

		return nil
	})

	return db, err
}
