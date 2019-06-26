package migrator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/mypackbot/internal/db"
)

func TestAutoMigrate(t *testing.T) {
	testDir := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com", "toby3d", "mypackbot", "test")
	oldPath := filepath.Join(testDir, "testing.old")
	newPath := filepath.Join(testDir, "testing.new")

	newDB, err := db.Open(newPath)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, newDB.Close())
		assert.NoError(t, os.Remove(newPath))
	}()

	assert.NoError(t, AutoMigrate(oldPath, newDB))
}
