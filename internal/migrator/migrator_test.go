package migrator

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/mypackbot/internal"
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
		// assert.NoError(t, os.Remove(newPath))
	}()

	bot, err := internal.New(filepath.Join("..", "..", "configs", "production.yaml"))
	if err != nil {
		log.Fatalln("ERROR:", err.Error())
	}

	assert.NoError(t, AutoMigrate(AutoMigrateConfig{
		FromPath: oldPath,
		ToConn:   bot.Store(),
		Bot:      bot.Bot(),
	}))
}
