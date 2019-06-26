package store

import (
	"os"
	"path/filepath"
	"testing"

	bolt "github.com/etcd-io/bbolt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/mypackbot/internal/db"
)

func newDB(t *testing.T) (*bolt.DB, func()) {
	t.Helper()

	path := filepath.Join(".", "testing.db")
	dataBase, err := bolt.Open(path, os.ModePerm, nil)
	assert.NoError(t, err)

	db.AutoMigrate(dataBase)
	return dataBase, func() {
		assert.NoError(t, dataBase.Close())
		assert.NoError(t, os.Remove(path))
	}
}
