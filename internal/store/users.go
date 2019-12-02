package store

import (
	"strconv"
	"time"

	bolt "github.com/etcd-io/bbolt"
	json "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	"gitlab.com/toby3d/mypackbot/internal/common"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"golang.org/x/xerrors"
)

type UsersStore struct {
	conn     *bolt.DB
	marshler json.API
	parser   fastjson.Parser
}

var (
	ErrUserExist = model.Error{
		Message: "User already exist",
	}

	ErrUserNotExist = model.Error{
		Message: "User not exist",
	}
)

func NewUsersStore(conn *bolt.DB, marshler json.API) *UsersStore {
	return &UsersStore{
		conn:     conn,
		marshler: marshler,
		parser:   fastjson.Parser{},
	}
}

func (store *UsersStore) Create(u *model.User) error {
	if store.Get(u.ID) != nil || store.GetByUserID(u.UserID) != nil {
		return ErrUserExist
	}

	now := time.Now().UTC().Unix()

	if u.CreatedAt <= 0 {
		u.CreatedAt = now
	}

	if u.UpdatedAt <= 0 {
		u.UpdatedAt = now
	}

	if u.LastSeen <= 0 {
		u.LastSeen = now
	}

	return store.conn.Update(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket(common.BucketUsers)

		if u.ID, err = bkt.NextSequence(); err != nil {
			return err
		}

		src, err := store.marshler.Marshal(u)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(strconv.FormatUint(u.ID, 10)), src)
	})
}

func (store *UsersStore) Get(id uint64) *model.User {
	u := new(model.User)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		return store.marshler.Unmarshal(
			tx.Bucket(common.BucketUsers).Get([]byte(strconv.FormatUint(id, 10))), u,
		)
	}); err != nil || u.ID == 0 {
		return nil
	}

	return u
}

func (store *UsersStore) GetByUserID(id int64) *model.User {
	u := new(model.User)

	if err := store.conn.View(func(tx *bolt.Tx) error {
		if err := tx.Bucket(common.BucketUsers).ForEach(func(key, val []byte) error {
			v, err := store.parser.ParseBytes(val)
			if err != nil {
				return err
			}

			if v.GetInt64("user_id") != id {
				return nil
			}

			if err = store.marshler.Unmarshal(val, u); err != nil {
				return err
			}

			return ErrForEachStop
		}); err != nil && !xerrors.Is(err, ErrForEachStop) {
			return err
		}

		return nil
	}); err != nil || u.ID == 0 {
		return nil
	}

	return u
}

func (store *UsersStore) Update(u *model.User) error {
	if store.Get(u.ID) == nil && store.GetByUserID(u.UserID) == nil {
		return store.Create(u)
	}

	if u.UpdatedAt <= 0 {
		u.UpdatedAt = time.Now().UTC().Unix()
	}

	src, err := store.marshler.Marshal(u)
	if err != nil {
		return err
	}

	return store.conn.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(common.BucketUsers).Put([]byte(strconv.FormatUint(u.ID, 10)), src)
	})
}

func (store *UsersStore) GetOrCreate(u *model.User) (user *model.User, err error) {
	if user = store.Get(u.ID); user != nil {
		return user, nil
	}

	if user = store.GetByUserID(u.UserID); user != nil {
		return user, nil
	}

	if err = store.Create(u); err != nil {
		return nil, err
	}

	return store.GetOrCreate(u)
}
