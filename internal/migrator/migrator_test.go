package migrator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/store"
)

func TestAutoMigrate(t *testing.T) {
	testDir := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com", "toby3d", "mypackbot", "test")
	oldPath := filepath.Join(testDir, "testing.old")
	newPath := filepath.Join(testDir, "testing.new")

	oldDB, err := bunt.Open(oldPath)
	assert.NoError(t, err)
	defer oldDB.Close()

	newDB, err := db.Open(newPath)
	assert.NoError(t, err)
	defer func() {
		newDB.Close()
		os.Remove(newPath)
	}()

	assert.NoError(t, AutoMigrate(AutoMigrateConfig{
		OldDB: oldDB,
		NewDB: store.NewStore(newDB),
	}))
}
