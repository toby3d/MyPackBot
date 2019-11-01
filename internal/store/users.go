package store

import (
	"errors"
	"strconv"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/json-iterator/go"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
)

type UsersStore struct{ conn *bolt.DB }

func NewUsersStore(conn *bolt.DB) *UsersStore { return &UsersStore{conn: conn} }

func (store *UsersStore) Create(u *model.User) error {
	if store.Get(u.ID) != nil {
		return errors.New("user already exists")
	}

	now := time.Now().UTC().Unix()

	if u.CreatedAt <= 0 {
		u.CreatedAt = now
	}

	if u.UpdatedAt <= 0 {
		u.UpdatedAt = now
	}

	src, err := json.ConfigFastest.Marshal(u)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsers).Put([]byte(strconv.Itoa(u.ID)), src)
	})
}

func (store *UsersStore) Get(uid int) *model.User {
	u := new(model.User)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		src := tx.Bucket(common.BucketUsers).Get([]byte(strconv.Itoa(uid)))

		return json.ConfigFastest.Unmarshal(src, u)
	}); err != nil {
		return nil
	}

	return u
}

func (store *UsersStore) Update(u *model.User) error {
	if store.Get(u.ID) == nil {
		return store.Create(u)
	}

	if u.UpdatedAt <= 0 {
		u.UpdatedAt = time.Now().UTC().Unix()
	}

	src, err := json.ConfigFastest.Marshal(u)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsers).Put([]byte(strconv.Itoa(u.ID)), src)
	})
}

func (store *UsersStore) Remove(uid int) error {
	if store.Get(uid) == nil {
		return errors.New("user already removed or not exists")
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(common.BucketUsers)

		if err := bkt.Delete([]byte(strconv.Itoa(uid))); err != nil {
			return err
		}

		bkt = tx.Bucket(common.BucketUsersStickers)

		return bkt.ForEach(func(key, val []byte) error {
			us := new(model.UserSticker)

			if err := json.Unmarshal(val, us); err != nil {
				return err
			}

			if us.UserID != uid {
				return nil
			}

			return bkt.Delete(key)
		})
	})
}

func (store *UsersStore) GetOrCreate(u *model.User) (*model.User, error) {
	if user := store.Get(u.ID); user != nil {
		return user, nil
	}

	if err := store.Create(u); err != nil {
		return nil, err
	}

	return store.Get(u.ID), nil
}
