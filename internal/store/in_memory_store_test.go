package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/mypackbot/internal/model"
)

func TestInMemoryUsersStore(t *testing.T) {
	u := model.User{
		ID:           42,
		CreatedAt:    time.Now().UTC().Unix(),
		UpdatedAt:    time.Now().UTC().Unix(),
		LanguageCode: "ru",
		LastSeen:     time.Now().UTC().Unix(),
	}

	store := NewInMemoryUsersStore()
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

func TestInMemoryStickersStore(t *testing.T) {
	s := model.Sticker{
		ID:         "abc",
		SetName:    "testing",
		IsAnimated: true,
		CreatedAt:  time.Now().UTC().Unix(),
	}

	store := NewInMemoryStickersStore()
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

func TestInMemoryStore(t *testing.T) {
	u := model.User{
		ID:           42,
		LanguageCode: "en",
	}
	s := model.Sticker{
		ID:      "abc",
		SetName: "testing",
		Emoji:   "👨‍💻💻",
	}

	store := NewInMemoryStore()
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

	stickers, count = store.GetStickersSet(&u, 0, 50, "wtf")
	assert.Equal(t, 0, count)
	assert.Len(t, stickers, 0)
	assert.Empty(t, stickers)

	assert.NoError(t, store.RemoveSticker(&u, &s))
}
