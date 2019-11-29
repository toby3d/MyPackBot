package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		db, err := Open(filepath.Join("/", "invalid", "path"))
		assert.Error(t, err)

		t.Run("automigrate", func(t *testing.T) {
			assert.Panics(t, func() { _ = AutoMigrate(db) })
		})
	})
	t.Run("valid", func(t *testing.T) {
		rootPath, err := os.Getwd()
		assert.NoError(t, err)

		testPath := filepath.Join(rootPath, "..", "..", "test", "testing.db")
		db, err := Open(testPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			assert.NoError(t, db.Close())
			assert.NoError(t, os.Remove(testPath))
		}()

		t.Run("automigrate", func(t *testing.T) {
			assert.NoError(t, AutoMigrate(db))
		})
	})
}
