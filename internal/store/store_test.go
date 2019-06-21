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
		if _, err = tx.CreateBucket([]byte("users")); err != nil {
			return err
		}

		if _, err = tx.CreateBucket([]byte("stickers")); err != nil {
			return err
		}

		if _, err = tx.CreateBucket([]byte("users_stickers")); err != nil {
			return err
		}

		return nil
	}))
	return db, func() {
		assert.NoError(t, db.Close())
		assert.NoError(t, os.Remove(path))
	}
}
