package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	bolt "github.com/etcd-io/bbolt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/model"
)

func initDB(t *testing.T) (*bolt.DB, func()) {
	rootPath, err := os.Getwd()
	assert.NoError(t, err)

	dbPath := filepath.Join(rootPath, "..", "..", "test", "testing.db")
	dataBase, err := db.Open(dbPath)

	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	return dataBase, func() {
		assert.NoError(t, dataBase.Close())
		assert.NoError(t, os.RemoveAll(dbPath))
	}
}

func TestUsersStore(t *testing.T) {
	u := model.User{
		ID:           42,
		CreatedAt:    time.Now().UTC().Unix(),
		UpdatedAt:    time.Now().UTC().Unix(),
		LanguageCode: "ru",
		LastSeen:     time.Now().UTC().Unix(),
	}

	dataBase, release := initDB(t)
	defer release()

	store := NewUsersStore(dataBase)
	assert.Nil(t, store.Get(u.ID))
	assert.NoError(t, store.Create(&u))
	assert.Equal(t, &u, store.Get(u.ID))
	assert.Error(t, store.Create(&u))
	assert.NoError(t, store.Remove(u.ID))
	assert.Error(t, store.Remove(u.ID))
	assert.NoError(t, store.Update(&u))
	assert.Equal(t, &u, store.Get(u.ID))
	assert.NoError(t, store.Update(&model.User{
		ID:           u.ID,
		LanguageCode: "en",
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    time.Now().UTC().Unix(),
		LastSeen:     time.Now().UTC().Unix(),
	}))
	assert.NotEqual(t, &u, store.Get(u.ID))
}

func TestStickersStore(t *testing.T) {
	s := model.Sticker{
		ID:         "abc",
		SetName:    "testing",
		IsAnimated: true,
		CreatedAt:  time.Now().UTC().Unix(),
	}

	dataBase, release := initDB(t)
	defer release()

	store := NewStickersStore(dataBase)
	assert.Nil(t, store.Get(s.ID))
	assert.NoError(t, store.Create(&s))
	assert.Equal(t, &s, store.Get(s.ID))
	assert.Error(t, store.Create(&s))
	assert.NoError(t, store.Remove(s.ID))
	assert.Error(t, store.Remove(s.ID))
	assert.NoError(t, store.Update(&s))
	assert.Equal(t, &s, store.Get(s.ID))
	assert.NoError(t, store.Update(&model.Sticker{
		ID:         s.ID,
		SetName:    "debug",
		IsAnimated: false,
		CreatedAt:  time.Now().UTC().Unix(),
	}))
	assert.NotEqual(t, &s, store.Get(s.ID))
}

func TestStore(t *testing.T) {
	u := model.User{
		ID:           42,
		LanguageCode: "en",
	}
	s := model.Sticker{
		ID:      "abc",
		SetName: "testing",
		Emoji:   "\u200d💻",
	}

	dataBase, release := initDB(t)
	defer release()

	store := NewStore(dataBase)
	assert.Error(t, store.RemoveSticker(&u, &s))

	assert.NoError(t, store.AddSticker(&u, &s))
	assert.Error(t, store.AddSticker(&u, &s))

	stickers, count := store.GetStickersList(&u, 0, 50, "")
	assert.Equal(t, 1, count)
	assert.Len(t, stickers, 1)
	assert.Contains(t, stickers, &s)

	stickers, count = store.GetStickersList(&u, 0, 50, "🐱")
	assert.Equal(t, 0, count)
	assert.Len(t, stickers, 0)
	assert.Empty(t, stickers)

	stickers, count = store.GetStickersSet(&u, 0, 50, s.SetName)
	assert.Equal(t, 1, count)
	assert.Len(t, stickers, 1)
	assert.Contains(t, stickers, &s)

	assert.NoError(t, store.AddStickersSet(&u, s.SetName))

	stickers, count = store.GetStickersSet(&u, 0, 50, "wtf")
	assert.Equal(t, 0, count)
	assert.Len(t, stickers, 0)
	assert.Empty(t, stickers)

	assert.NoError(t, store.RemoveSticker(&u, &s))
}
