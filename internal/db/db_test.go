package db

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

var (
	path   = filepath.Join(os.Getenv("GOPATH"), "src", "gitlab.com", "toby3d", "mypackbot", "test")
	dbPath = filepath.Join(path, "testing.db")

	user = User{
		ID:       42,
		Language: language.English,
	}
	set     = Set{ID: "set"}
	sticker = Sticker{ID: "sticker"}
)

func TestMain(m *testing.M) {
	set.Stickers = []*Sticker{&sticker}
	set.User = &user
	user.Sets = []*Set{&set}
	sticker.User = &user
	sticker.Set = &set

	os.Exit(m.Run())
}

func TestOpen(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open("/not/valid/path")
		assert.Error(t, err)
		assert.Nil(t, db)
		t.Run("close", func(t *testing.T) {
			assert.Error(t, db.Close())
			os.Remove(dbPath)
		})
	})

	t.Run("valid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		t.Run("close", func(t *testing.T) {
			assert.NoError(t, db.Close())
			os.Remove(dbPath)
		})
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		err := db.CreateUser(&user)
		if assert.Error(t, err) {
			assert.EqualError(t, ErrDatabaseClosed, err.Error())
		}
	})

	t.Run("valid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateUser(&user))
	})
}

func TestCreateSet(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		err := db.CreateSet(&set)
		if assert.Error(t, err) {
			assert.EqualError(t, ErrDatabaseClosed, err.Error())
		}
	})

	t.Run("valid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateSet(&set))
	})
}

func TestCreateSticker(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		err := db.CreateSticker(&sticker)
		if assert.Error(t, err) {
			assert.EqualError(t, ErrDatabaseClosed, err.Error())
		}
	})

	t.Run("valid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateSticker(&sticker))
	})
}

func TestGetUser(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		u, err := db.GetUser(&user)
		if assert.Error(t, err) {
			assert.EqualError(t, ErrDatabaseClosed, err.Error())
		}
		assert.Nil(t, u)
	})

	t.Run("valid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		if assert.NoError(t, db.CreateUser(&user)) {
			u, err := db.GetUser(&user)
			assert.NoError(t, err)
			if assert.NotNil(t, u) {
				assert.Equal(t, &user, u)
			}
		}
	})
}

func TestGetSet(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		s, err := db.GetSet(&set)
		if assert.Error(t, err) {
			assert.EqualError(t, ErrDatabaseClosed, err.Error())
		}
		assert.Nil(t, s)
	})
	t.Run("valid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		if assert.NoError(t, db.CreateSet(&set)) {
			s, err := db.GetSet(&set)
			assert.NoError(t, err)
			if assert.NotNil(t, s) {
				assert.Equal(t, &set, s)
			}
		}
	})
}

func TestGetSticker(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		s, err := db.GetSticker(&sticker)
		if assert.Error(t, err) {
			assert.EqualError(t, ErrDatabaseClosed, err.Error())
		}
		assert.Nil(t, s)
	})

	t.Run("valid", func(t *testing.T) {
		os.Remove(dbPath)

		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		if assert.NoError(t, db.CreateSticker(&sticker)) {
			s, err := db.GetSticker(&sticker)
			assert.NoError(t, err)
			if assert.NotNil(t, s) {
				assert.Equal(t, &sticker, s)
			}
		}
	})
}

func TestUpdateUser(t *testing.T) {
	defer func() {
		user.Language = language.English
		user.Hits = 0
	}()

	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		assert.Error(t, db.UpdateUser(&user))
	})

	t.Run("valid", func(t *testing.T) {
		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateUser(&user))

		user.Language = language.Russian
		user.Hits = rand.Intn(100)
		assert.NoError(t, db.UpdateUser(&user))

		u, err := db.GetUser(&user)
		if assert.NoError(t, err) {
			assert.Equal(t, user.Hits, u.Hits)
			assert.Equal(t, user.Language, u.Language)
		}
	})
}

func TestUpdateSet(t *testing.T) {
	defer func() {
		set.Hits = 0
		set.IsFavorite = false
	}()

	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		assert.Error(t, db.UpdateSet(&set))
	})

	t.Run("valid", func(t *testing.T) {
		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateSet(&set))

		set.Hits = rand.Intn(100)
		set.IsFavorite = true
		assert.NoError(t, db.UpdateSet(&set))

		s, err := db.GetSet(&set)
		if assert.NoError(t, err) {
			assert.Equal(t, set.Hits, s.Hits)
			assert.Equal(t, set.IsFavorite, s.IsFavorite)
		}
	})
}

func TestUpdateSticker(t *testing.T) {
	defer func() {
		sticker.Hits = 0
		sticker.IsFavorite = false
	}()

	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		assert.Error(t, db.UpdateSticker(&sticker))
	})

	t.Run("valid", func(t *testing.T) {
		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateSticker(&sticker))

		sticker.Hits = rand.Intn(100)
		sticker.IsFavorite = true
		assert.NoError(t, db.UpdateSticker(&sticker))

		s, err := db.GetSticker(&sticker)
		if assert.NoError(t, err) {
			assert.Equal(t, sticker.Hits, s.Hits)
			assert.Equal(t, sticker.IsFavorite, s.IsFavorite)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		assert.Error(t, db.DeleteUser(&user))
	})

	t.Run("valid", func(t *testing.T) {
		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateUser(&user))
		assert.NoError(t, db.DeleteUser(&user))
	})
}

func TestDeleteSet(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		assert.Error(t, db.DeleteSet(&set))
	})

	t.Run("valid", func(t *testing.T) {
		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateSet(&set))
		assert.NoError(t, db.DeleteSet(&set))
	})
}

func TestDeleteSticker(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		os.Remove(dbPath)
		defer os.Remove(dbPath)

		var db DB
		assert.Error(t, db.DeleteSticker(&sticker))
	})

	t.Run("valid", func(t *testing.T) {
		db, err := Open(dbPath)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer func() {
			db.Close()
			os.Remove(dbPath)
		}()

		assert.NoError(t, db.CreateSticker(&sticker))
		assert.NoError(t, db.DeleteSticker(&sticker))
	})
}
