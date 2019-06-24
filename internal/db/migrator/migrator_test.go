package migrator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/db"
)

func TestAutoMigrate(t *testing.T) {
	testDir := filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com", "toby3d", "mypackbot", "test")
	oldPath := filepath.Join(testDir, "testing.old")
	newPath := filepath.Join(testDir, "testing.new")

	oldDB, err := bunt.Open(oldPath)
	assert.NoError(t, err)
	assert.NoError(t, oldDB.Update(func(tx *bunt.Tx) (err error) {
		if _, _, err = tx.Set("user:123:state", "idle", nil); err != nil {
			return err
		}
		if _, _, err = tx.Set("user:123:set:?:sticker:abc", "👌", nil); err != nil {
			return err
		}
		_, _, err = tx.Set("user:123:set:testing:sticker:cba", "👍", nil)
		return err
	}))
	assert.NoError(t, oldDB.Close())

	newDB, err := db.Open(newPath)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, newDB.Close())
		assert.NoError(t, os.Remove(newPath))
	}()

	assert.NoError(t, AutoMigrate(oldPath, newDB))
}
