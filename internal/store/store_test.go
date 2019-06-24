package store

import (
	"os"
	"path/filepath"
	"testing"

	bolt "github.com/etcd-io/bbolt"
	"github.com/stretchr/testify/assert"
)

func newDB(t *testing.T) (*bolt.DB, func()) {
	t.Helper()

	path := filepath.Join(".", "testing.db")
	db, err := bolt.Open(path, os.ModePerm, nil)
	assert.NoError(t, err)

	assert.NoError(t, db.Update(func(tx *bolt.Tx) (err error) {
		if _, err = tx.CreateBucket(bktUsers); err != nil {
			return err
		}

		if _, err = tx.CreateBucket(bktStickers); err != nil {
			return err
		}

		if _, err = tx.CreateBucket(bktUsersStickers); err != nil {
			return err
		}

		if _, err = tx.CreateBucket(bktChannels); err != nil {
			return err
		}

		return nil
	}))
	return db, func() {
		assert.NoError(t, db.Close())
		assert.NoError(t, os.Remove(path))
	}
}
