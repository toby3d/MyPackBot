package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/mypackbot/internal/models"
)

func TestGetStickerByID(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewStickerStore(db)
	s := models.Sticker{
		ID:      "abc",
		Emoji:   "👍",
		SetName: "testing",
	}
	assert.NoError(t, store.Create(&s))

	t.Run("invalid", func(t *testing.T) {
		sticker, err := store.GetByID("cba")
		assert.NoError(t, err)
		assert.Empty(t, sticker)
	})
	t.Run("valid", func(t *testing.T) {
		sticker, err := store.GetByID("abc")
		assert.NoError(t, err)
		assert.Equal(t, &s, sticker)
	})
}

func TestStickersGetByUserID(t *testing.T) {
	db, release := newDB(t)
	defer release()

	stickerStore := NewStickerStore(db)
	stickers := []models.Sticker{
		models.Sticker{
			ID:      "cba",
			Emoji:   "👌",
			SetName: "test",
		},
		models.Sticker{
			ID:      "abc",
			Emoji:   "👍",
			SetName: "testing",
		},
		models.Sticker{
			ID:      "123",
			Emoji:   "😺",
			SetName: "test",
		},
		models.Sticker{
			ID:      "321",
			Emoji:   "👌",
			SetName: "testing",
		},
	}
	usersStickers := []models.UsersStickers{
		models.UsersStickers{
			UserID:    42,
			StickerID: "cba",
		},
		models.UsersStickers{
			UserID:    42,
			StickerID: "123",
		},
		models.UsersStickers{
			UserID:    42,
			StickerID: "321",
		},
	}
	for _, s := range stickers {
		s := s
		assert.NoError(t, stickerStore.Create(&s))
	}
	userStore := NewUserStore(db)
	for _, us := range usersStickers {
		us := us
		assert.NoError(t, userStore.AddSticker(us.UserID, us.StickerID))
	}

	var empty []models.Sticker
	for _, tc := range []struct {
		info      string
		query     int
		offset    int
		limit     int
		expResult []models.Sticker
		expCount  int
	}{{
		info:      "all of 42",
		offset:    -1,
		limit:     -1,
		query:     42,
		expResult: []models.Sticker{stickers[0], stickers[2], stickers[3]},
		expCount:  3,
	}, {
		info:      "first of 42",
		offset:    -1,
		limit:     1,
		query:     42,
		expResult: []models.Sticker{stickers[3]},
		expCount:  3,
	}, {
		info:      "1-2 of 42",
		offset:    -1,
		limit:     2,
		query:     42,
		expResult: []models.Sticker{stickers[2], stickers[3]},
		expCount:  3,
	}, {
		info:      "2-3 of 42",
		offset:    1,
		limit:     -1,
		query:     42,
		expResult: []models.Sticker{stickers[0], stickers[2]},
		expCount:  3,
	}, {
		info:      "none 24",
		offset:    -1,
		limit:     -1,
		query:     24,
		expResult: empty,
		expCount:  0,
	}} {
		tc := tc
		t.Run(tc.info, func(t *testing.T) {
			list, count, err := stickerStore.GetByUserID(tc.query, tc.offset, tc.limit)
			assert.NoError(t, err)
			assert.Equal(t, tc.expCount, count)
			for _, r := range tc.expResult {
				r := r
				assert.Contains(t, list, r)
			}
		})
	}
}

func TestCreateSticker(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewStickerStore(db)
	s := models.Sticker{
		ID:      "abc",
		Emoji:   "👍",
		SetName: "testing",
	}

	t.Run("invalid", func(t *testing.T) {
		assert.Error(t, store.Create(&models.Sticker{}))
	})
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, store.Create(&s))
	})
}

func TestUpdateSticker(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewStickerStore(db)
	s := models.Sticker{
		ID:      "abc",
		Emoji:   "👍",
		SetName: "testing",
	}
	assert.NoError(t, store.Create(&s))

	t.Run("invalid", func(t *testing.T) {
		assert.Error(t, store.Update(&models.Sticker{}))
	})
	t.Run("valid", func(t *testing.T) {
		s2 := models.Sticker{
			ID:      "abc",
			Emoji:   "👌",
			SetName: "testing",
		}
		assert.NoError(t, store.Update(&s2))
		assert.NotEqual(t, s2, s)
	})
}

func TestDeleteSticker(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewStickerStore(db)
	s := models.Sticker{
		ID:      "abc",
		Emoji:   "👍",
		SetName: "testing",
	}
	assert.NoError(t, store.Create(&s))

	t.Run("invalid", func(t *testing.T) {
		assert.Error(t, store.Update(&models.Sticker{}))
	})
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, store.Delete(&s))
	})
}

func TestStickersList(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewStickerStore(db)
	stickers := []models.Sticker{
		models.Sticker{
			ID:      "cba",
			Emoji:   "👌",
			SetName: "testing",
		},
		models.Sticker{
			ID:      "abc",
			Emoji:   "👍",
			SetName: "testing",
		},
		models.Sticker{
			ID:      "123",
			Emoji:   "👋",
			SetName: "testing",
		},
	}
	for _, s := range stickers {
		s := s
		assert.NoError(t, store.Create(&s))
	}

	var empty []models.Sticker
	for _, tc := range []struct {
		info      string
		offset    int
		limit     int
		expResult []models.Sticker
	}{{
		info:      "get all",
		offset:    -1,
		limit:     -1,
		expResult: stickers,
	}, {
		info:      "2-3",
		offset:    1,
		limit:     -1,
		expResult: stickers[1:],
	}, {
		info:      "3",
		offset:    2,
		limit:     -1,
		expResult: stickers[2:],
	}, {
		info:      "nil",
		offset:    3,
		limit:     -1,
		expResult: empty,
	}, {
		info:      "1",
		offset:    0,
		limit:     1,
		expResult: stickers[:1],
	}, {
		info:      "1-2",
		offset:    0,
		limit:     2,
		expResult: stickers[:2],
	}, {
		info:      "2",
		offset:    1,
		limit:     1,
		expResult: stickers[1:2],
	}, {
		info:      "2-3",
		offset:    1,
		limit:     2,
		expResult: stickers[1:],
	}} {
		tc := tc
		t.Run(tc.info, func(t *testing.T) {
			list, count, err := store.List(tc.offset, tc.limit)
			assert.NoError(t, err)
			assert.Equal(t, len(stickers), count)
			for _, r := range tc.expResult {
				r := r
				assert.Contains(t, list, r)
			}
		})
	}
}

func TestStickersListByEmoji(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewStickerStore(db)
	stickers := []models.Sticker{
		models.Sticker{
			ID:      "cba",
			Emoji:   "👌",
			SetName: "testing",
		},
		models.Sticker{
			ID:      "abc",
			Emoji:   "👍",
			SetName: "testing",
		},
		models.Sticker{
			ID:      "123",
			Emoji:   "👌",
			SetName: "testing",
		},
	}
	for _, s := range stickers {
		s := s
		assert.NoError(t, store.Create(&s))
	}

	var empty []models.Sticker
	for _, tc := range []struct {
		info      string
		query     string
		offset    int
		limit     int
		expResult []models.Sticker
		expCount  int
	}{{
		info:      "all 👌",
		offset:    -1,
		limit:     -1,
		query:     "👌",
		expResult: []models.Sticker{stickers[0], stickers[2]},
		expCount:  2,
	}, {
		info:      "all 👍",
		offset:    -1,
		limit:     -1,
		query:     "👍",
		expResult: []models.Sticker{stickers[1]},
		expCount:  1,
	}, {
		info:      "first 👌",
		limit:     1,
		offset:    -1,
		query:     "👌",
		expResult: []models.Sticker{stickers[0]},
		expCount:  2,
	}, {
		info:      "second 👍",
		limit:     1,
		offset:    1,
		query:     "👍",
		expResult: empty,
		expCount:  1,
	}} {
		tc := tc
		t.Run(tc.info, func(t *testing.T) {
			list, count, err := store.ListByEmoji(tc.query, tc.offset, tc.limit)
			assert.NoError(t, err)
			assert.Equal(t, tc.expCount, count)
			for _, r := range tc.expResult {
				r := r
				assert.Contains(t, list, r)
			}
		})
	}
}

func TestStickersGetSet(t *testing.T) {
	db, release := newDB(t)
	defer release()

	store := NewStickerStore(db)
	stickers := []models.Sticker{
		models.Sticker{
			ID:      "cba",
			Emoji:   "👌",
			SetName: "test",
		},
		models.Sticker{
			ID:      "abc",
			Emoji:   "👍",
			SetName: "testing",
		},
		models.Sticker{
			ID:      "123",
			Emoji:   "😺",
			SetName: "test",
		},
		models.Sticker{
			ID:      "321",
			Emoji:   "👌",
			SetName: "testing",
		},
	}
	for _, s := range stickers {
		s := s
		assert.NoError(t, store.Create(&s))
	}

	var empty []models.Sticker
	for _, tc := range []struct {
		info      string
		query     string
		offset    int
		limit     int
		expResult []models.Sticker
		expCount  int
	}{{
		info:      "all test",
		offset:    -1,
		limit:     -1,
		query:     "test",
		expResult: []models.Sticker{stickers[0], stickers[2]},
		expCount:  2,
	}, {
		info:      "all testing",
		offset:    -1,
		limit:     -1,
		query:     "testing",
		expResult: []models.Sticker{stickers[1], stickers[3]},
		expCount:  2,
	}, {
		info:      "empty wtf",
		offset:    -1,
		limit:     -1,
		query:     "wtf",
		expResult: empty,
		expCount:  0,
	}, {
		info:      "first test",
		offset:    -1,
		limit:     1,
		query:     "test",
		expResult: []models.Sticker{stickers[0]},
		expCount:  2,
	}, {
		info:      "second testing",
		offset:    1,
		limit:     -1,
		query:     "testing",
		expResult: []models.Sticker{stickers[3]},
		expCount:  2,
	}} {
		tc := tc
		t.Run(tc.info, func(t *testing.T) {
			list, count, err := store.GetSet(tc.query, tc.offset, tc.limit)
			assert.NoError(t, err)
			assert.Equal(t, tc.expCount, count)
			for _, r := range tc.expResult {
				r := r
				assert.Contains(t, list, r)
			}
		})
	}
}
