package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/mypackbot/internal/models"
)

func TestGetUserByID(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewUserStore(db)
	u := models.User{
		ID:       42,
		Language: "ru",
	}
	assert.NoError(t, store.Create(&u))

	t.Run("invalid", func(t *testing.T) {
		user, err := store.GetByID(24)
		assert.NoError(t, err)
		assert.Empty(t, user)
	})
	t.Run("valid", func(t *testing.T) {
		user, err := store.GetByID(42)
		assert.NoError(t, err)
		assert.Equal(t, &u, user)
	})
}

func TestCreateUser(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewUserStore(db)
	u := models.User{
		ID:       42,
		Language: "ru",
	}

	assert.NoError(t, store.Create(&u))
}

func TestUpdateUser(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewUserStore(db)
	s := models.User{
		ID:       42,
		Language: "ru",
	}
	assert.NoError(t, store.Create(&s))

	s2 := models.User{
		ID:       42,
		Language: "en",
	}
	assert.NoError(t, store.Update(&s2))
	assert.NotEqual(t, s2, s)
}

func TestUserAddSticker(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewUserStore(db)
	s := models.User{
		ID:       42,
		Language: "ru",
	}
	assert.NoError(t, store.Create(&s))

	assert.NoError(t, store.AddSticker(s.ID, "abc"))
}

func TestUserDeleteSticker(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewUserStore(db)
	s := models.User{
		ID:       42,
		Language: "ru",
	}
	assert.NoError(t, store.Create(&s))

	assert.NoError(t, store.AddSticker(s.ID, "abc"))
	assert.NoError(t, store.DeleteSticker(s.ID, "abc"))
}
